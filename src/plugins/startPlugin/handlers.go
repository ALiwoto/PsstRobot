package startPlugin

import (
	"strconv"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/src/core"
	"github.com/AnimeKaizoku/PsstRobot/src/core/utils"
	wv "github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	"github.com/AnimeKaizoku/PsstRobot/src/database/usersDatabase"
	"github.com/AnimeKaizoku/PsstRobot/src/database/whisperDatabase"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func startHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	// checking for a space in text message is most probably one of
	// the best ways in reducing checkers in un-expected cases.
	if strings.Contains(msg.Text, " ") {
		// and example of a start handler with a cb param can be this:
		// "/start 1gs9860=17uurds"
		if utils.IsStartedForWhisper(ctx.EffectiveMessage.Text) {
			return sendWhisperText(bot, ctx)
		}

		// user might have been redirected from other inline
		// sections of the bot (such as "too long" error).
		myStrs := strings.Split(msg.Text, " ")
		if len(myStrs) > 1 {
			switch myStrs[1] {
			case wv.StartDataCreate:
				return wv.CreateWhisperHandler(bot, ctx)
			case wv.HelpDataInline:
				return helpHandler(bot, ctx)
			}
		}
	}

	return normalStartHandler(bot, ctx)
}

func privacyHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage
	args := ssg.SplitN(message.Text, 2, " ")
	enabled := usersDatabase.HasPrivacy(user)
	preString := func() string {
		if enabled {
			return "enabled"
		}
		return "disabled"
	}

	var md mdparser.WMarkDown
	shouldEnable, ok := utils.ArgsToBool(args)
	if ok {
		if shouldEnable == enabled {
			md = mdparser.GetNormal("Privacy mode is already " + preString() + ".")
		} else {
			enabled = shouldEnable
			usersDatabase.ChangePrivacy(user, enabled)
			md = mdparser.GetNormal("Privacy mode has been " + preString() + " successfully.")
		}
	} else {
		md = mdparser.GetBold("· Usage:")
		md.Normal("\n    /privacy [on|off]\n\n" +
			"Currently privacy mode is " + preString() +
			" for every whisper you try to read.",
		)
	}

	_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
	})
	return ext.EndGroups
}

func helpHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage

	md := mdparser.GetNormal("Hey ")
	md.Mention(user.FirstName, user.Id)
	md.Normal("!\nI'm " + bot.Username + ", a whisper bot with extra" +
		" features!\n")
	md.Normal("Here is some of the commands you can use:\n\n")

	md.Bold("• ").Mono("/create").Bold(": \n\t\t")
	md.Normal("Prepares an advanced whisper message, all types of medias and ")
	md.Normal("media-albums are supported.\n\n")
	md.Bold("• ").Mono("/privacy").Bold(": \n\t\t")
	md.Normal("Displays your current privacy mode. You can disable/enable it ")
	md.Normal("using this same command. If your privacy is enabled, ")
	md.Normal("you won't cause the whisper message to be edited to ")
	md.Normal("\"").Mention(user.FirstName, user.Id).Normal(" read the whisper\".\n\n")
	md.Bold("• ").Mono("/notifications").Bold(": \n\t\t")
	md.Normal("Displays your current notifications preferences using a panel.")
	md.Normal(" You can change your whisper-notifications preferences using ")
	md.Normal("using buttons of that panel.")

	_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             core.MarkdownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}

func normalStartHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage

	md := mdparser.GetNormal("Hello ")
	md.Mention(user.FirstName, user.Id)
	md.Normal("!\nThis bot allows you to send whisper messages")
	md.Normal(" to other users.\n")
	md.Normal("You don't need to add this bot in any group,")
	md.Normal(" it works in inline mode only and can send advanced")
	md.Normal(" whispers (long whispers and whispers with media).")

	_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             core.MarkdownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}

func sendWhisperText(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	uniqueId := utils.GetWhisperUniqueId(message)
	w := whisperDatabase.GetWhisper(uniqueId)
	if w == nil || !w.CanRead(user) {
		return normalStartHandler(bot, ctx)
	}

	senderInfo, err := bot.GetChat(w.Sender)
	var senderMd mdparser.WMarkDown
	if err == nil && senderInfo != nil {
		senderMd = mdparser.GetUserMention(senderInfo.FirstName, senderInfo.Id)
	} else {
		senderMd = mdparser.GetMono(strconv.FormatInt(senderInfo.Id, 10))
	}

	var targetMd mdparser.WMarkDown
	if w.Recipient != 0 && w.RecipientUsername == "" {
		targetInfo, err := bot.GetChat(w.Recipient)
		if err == nil && targetInfo != nil {
			targetMd = mdparser.GetUserMention(targetInfo.FirstName, targetInfo.Id)
		} else {
			targetMd = mdparser.GetMono(strconv.FormatInt(w.Recipient, 10))
		}
	} else {
		if w.RecipientUsername == "" {
			// everyone
			targetMd = mdparser.GetBold("everyone")
		} else {
			targetMd = mdparser.GetNormal(w.RecipientUsername)
		}
	}

	text := mdparser.GetBold("🔹This whisper has been sent from ").AppendThis(senderMd)
	text.Bold(" to ").AppendThis(targetMd)
	text.Bold(" at ").Mono(w.CreatedAt.UTC().Format("2006-01-02 15:04:05"))
	text.Bold(":\n\n")
	if w.Type == whisperDatabase.WhisperTypePlainText {
		text.Normal(w.Text)
	}

	_, _ = message.Reply(bot, text.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
	})

	if w.Type != whisperDatabase.WhisperTypePlainText {
		if w.IsMediaGroup() {
			_, _ = bot.SendMediaGroup(user.Id, w.GetMediaGroup(), nil)
		} else {
			switch w.Type {
			//WhisperTypePhoto
			//WhisperTypeVideo
			//WhisperTypeAudio
			//WhisperTypeVoice
			//WhisperTypeSticker
			//WhisperTypeDocument
			//WhisperTypeVideoNote
			//WhisperTypeAnimation
			//WhisperTypeDice
			case whisperDatabase.WhisperTypePhoto:
				_, _ = bot.SendPhoto(user.Id, w.FileId, &gotgbot.SendPhotoOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeVideo:
				_, _ = bot.SendVideo(user.Id, w.FileId, &gotgbot.SendVideoOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeAudio:
				_, _ = bot.SendAudio(user.Id, w.FileId, &gotgbot.SendAudioOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeVoice:
				_, _ = bot.SendVoice(user.Id, w.FileId, &gotgbot.SendVoiceOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeSticker:
				_, _ = bot.SendSticker(user.Id, w.FileId, &gotgbot.SendStickerOpts{
					//Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeDocument:
				_, _ = bot.SendDocument(user.Id, w.FileId, &gotgbot.SendDocumentOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeVideoNote:
				_, _ = bot.SendVideoNote(user.Id, w.FileId, &gotgbot.SendVideoNoteOpts{
					//Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeAnimation:
				_, _ = bot.SendAnimation(user.Id, w.FileId, &gotgbot.SendAnimationOpts{
					Caption: w.Text,
				})
			case whisperDatabase.WhisperTypeDice:
				_, _ = bot.SendDice(user.Id, &gotgbot.SendDiceOpts{
					Emoji: w.Text,
				})
			}
		}

	}

	if w.ShouldMarkAsRead(user) && !usersDatabase.HasPrivacy(user) {
		md := mdparser.GetUserMention(user.FirstName, user.Id)
		md.Normal(" read the whisper")
		_, _, _ = bot.EditMessageText(md.ToString(), &gotgbot.EditMessageTextOpts{
			InlineMessageId:       w.InlineMessageId,
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})

		whisperDatabase.RemoveWhisper(w)
		usersDatabase.SaveInHistory(w.Sender, user)
	}
	return ext.EndGroups
}
