package startPlugin

import (
	"strconv"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/usersDatabase"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func startHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	// and example of a start handler with a cb param can be this:
	// "/start 1gs9860=17uurds"
	if utils.IsStartedForWhisper(ctx.EffectiveMessage.Text) {
		return sendWhisperText(bot, ctx)
	}

	return normalStartHandler(bot, ctx)
}

func normalStartHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage

	md := mdparser.GetNormal("Hello ")
	md.AppendMentionThis(user.FirstName, user.Id)
	md.AppendNormalThis("!\nThis bot allows you to send whisper messages")
	md.AppendNormalThis(" to other users.\n")
	md.AppendNormalThis("You don't need to add this bot in any group,")
	md.AppendNormalThis(" it works in inline mode only and can send advanced")
	md.AppendNormalThis(" whispers (longs whispers and whispers with media).")

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
	if w == nil {
		// TODO:
		// display original welcome message which is on
		// normal /start command :0
		return normalStartHandler(bot, ctx)
	}

	if !w.CanRead(user) {
		// TODO:
		// display original welcome message which is on
		// normal /start command :0
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
			targetMd = mdparser.GetMono(strconv.FormatInt(targetInfo.Id, 10))
		}
	} else {
		targetMd = mdparser.GetNormal(w.RecipientUsername)
	}

	text := mdparser.GetBold("🔹This whisper has been sent from ").AppendThis(senderMd)
	text.AppendBoldThis(" to ").AppendThis(targetMd).AppendBoldThis(":\n\n")
	text.AppendNormalThis(w.Text)

	_, _ = message.Reply(bot, text.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
	})

	if w.ShouldMarkAsRead(user) && usersDatabase.HasPrivacy(user) {
		md := mdparser.GetUserMention(user.FirstName, user.Id)
		md.AppendNormalThis(" read the whisper")
		_, _ = bot.EditMessageText(md.ToString(), &gotgbot.EditMessageTextOpts{
			InlineMessageId:       w.InlineMessageId,
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})

		whisperDatabase.RemoveWhisper(w)
		usersDatabase.SaveInHistory(w.Sender, user)
	}
	return ext.EndGroups
}
