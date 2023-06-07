package helpPlugin

import (
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/src/core"
	"github.com/AnimeKaizoku/PsstRobot/src/database/usersDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

//---------------------------------------------------------

func helpCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return cq.Data == helpCommand
}

func helpHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage

	md := mdparser.GetNormal("Hey ")
	md.Mention(user.FirstName, user.Id)
	md.Normal("!\nI'm " + bot.Username + ", a whisper bot with extra" +
		" features!\n")
	md.Normal("You can use this in inline mode, just type:\n")
	md.Mono("@" + bot.Username + " @username your whisper text here!\n")
	md.Normal("(you can also use numeric user-id instead of a username).\n\n")
	md.Normal("Here is some of the commands you can use in my pm:\n")

	md.Bold("• ").Mono("/create").Bold(": \n\t\t")
	md.Normal("Prepares an advanced whisper message, all types of medias and ")
	md.Normal("media-albums are supported.\n\n")
	md.Bold("• ").Mono("/list").Bold(": \n\t\t")
	md.Normal("Lists all of your current pending whisper messages.\n\n")
	md.Bold("• ").Mono("/privacy").Bold(": \n\t\t")
	md.Normal("Displays your current privacy mode. You can disable/enable it ")
	md.Normal("using this same command. If your privacy is enabled, ")
	md.Normal("you won't cause the whisper message to be edited to ")
	md.Normal("\"").Mention(user.FirstName, user.Id).Normal(" read the whisper\".\n\n")
	md.Bold("• ").Mono("/notifications").Bold(": \n\t\t")
	md.Normal("Displays your current notifications preferences using a panel.")
	md.Normal(" You can change your whisper-notifications preferences using ")
	md.Normal("using buttons of that panel.")

	if message.From.Id == bot.Id {
		_, _, _ = message.EditText(bot, md.ToString(), &gotgbot.EditMessageTextOpts{
			ReplyMarkup:           *getMainMenuHelpButtons(),
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})
	} else {
		_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
			ReplyMarkup:           getMainMenuHelpButtons(),
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})
	}

	return ext.EndGroups
}

//---------------------------------------------------------

func userWhisperHistoryCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return cq.Data == userWhisperHistoryData
}

func userWhisperHistoryResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage

	txt := mdparser.GetBold("🗒User whisper history")
	txt.Normal("\n" + bot.Username + " keeps track of the three most recent users you've")
	txt.Normal(" sent whispers to, making it effortless for you to choose")
	txt.Normal(" them when sending new whispers. If you prefer, you can")
	txt.Normal(" also clear the list of recent users or disable this feature entirely.")
	txt.Normal("\n\nPlease do notice that disabling/enabling this feature does not")
	txt.Normal(" affect the whispers you have sent in past.")

	if message.From.Id == bot.Id {
		_, _, _ = message.EditText(bot, txt.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode:   gotgbot.ParseModeMarkdownV2,
			ReplyMarkup: *getUserWhisperHistoryButtons(),
		})
	} else {
		_, _ = message.Reply(bot, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeMarkdownV2,
			ReplyMarkup: getUserWhisperHistoryButtons(),
		})
	}

	return ext.EndGroups
}

func clearUserHistoryCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return cq.Data == clearUserHistoryData
}

func clearUserHistoryResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage

	var txt mdparser.WMarkDown

	if usersDatabase.ClearUserWhisperHistory(message.From.Id) {
		txt = mdparser.GetBold("User whisper history cleared successfully!")
	} else {
		txt = mdparser.GetBold("No user history to clear!")
	}

	if message.From.Id == bot.Id {
		_, _, _ = message.EditText(bot, txt.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
	} else {
		_, _ = message.Reply(bot, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
	}

	return ext.EndGroups
}

func disableUserHistoryCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return cq.Data == clearUserHistoryData
}

func disableUserHistoryResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage

	usersDatabase.DisableUserWhisperHistory(ctx.EffectiveUser)

	txt := mdparser.GetBold("User whisper history has been disabled.")
	if message.From.Id == bot.Id {
		_, _, _ = message.EditText(bot, txt.ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
	} else {
		_, _ = message.Reply(bot, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
	}

	return ext.EndGroups
}

//---------------------------------------------------------
