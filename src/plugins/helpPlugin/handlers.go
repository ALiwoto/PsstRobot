package helpPlugin

import (
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/src/core"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

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
	md.Normal("Lists all of your current pending whisper messages.")
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
		ReplyMarkup:           getMainMenuHelpButtons(),
		ParseMode:             core.MarkdownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}
