package whisperPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func answerForLongAdvanced(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.InlineQuery.Answer(bot, nil, &gotgbot.AnswerInlineQueryOpts{
		SwitchPmText:      "Too long! Use an advanced whisper!",
		SwitchPmParameter: "create",
	})
	return ext.EndGroups
}

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	sendWhisperIq := handlers.NewInlineQuery(sendwhisperFilter, sendWhisperResponse)
	chosenWhisperIq := handlers.NewChosenInlineResult(chosenWhisperFilter, chosenWhisperResponse)
	showWishperCb := handlers.NewCallback(showWhisperCallBackQuery, showWhisperResponse)

	d.AddHandler(chosenWhisperIq)
	d.AddHandler(sendWhisperIq)
	d.AddHandler(showWishperCb)
}
