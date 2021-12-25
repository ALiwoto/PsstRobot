package startPlugin

import (
	"strconv"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
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
		return ext.EndGroups
	}

	if !w.CanRead(user) {
		// TODO:
		// display original welcome message which is on
		// normal /start command :0
		return ext.EndGroups
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

	bot.SendSticker(0, "", &gotgbot.SendStickerOpts{})

	if w.InlineMessageId != "" && w.Sender != user.Id {
		go whisperDatabase.RemoveWhisper(w)

		md := mdparser.GetUserMention(user.FirstName, user.Id)
		md.AppendNormalThis(" read the whisper")
		_, _ = bot.EditMessageText(md.ToString(), &gotgbot.EditMessageTextOpts{
			InlineMessageId:       w.InlineMessageId,
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})
	}
	return ext.EndGroups
}
