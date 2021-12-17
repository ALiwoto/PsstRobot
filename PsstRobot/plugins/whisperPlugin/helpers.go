package whisperPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	sendWhisperIq := handlers.NewInlineQuery(sendwhisperFilter, sendWhisperResponse)
	chosenWhisperIq := handlers.NewChosenInlineResult(chosenWhisperFilter, chosenWhisperResponse)
	showWishperCb := handlers.NewCallback(showWhisperCallBackQuery, showWhisperResponse)

	d.AddHandler(chosenWhisperIq)
	d.AddHandler(sendWhisperIq)
	d.AddHandler(showWishperCb)
}
