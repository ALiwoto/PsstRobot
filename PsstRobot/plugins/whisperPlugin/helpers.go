package whisperPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadHandlers(d *ext.Dispatcher, t []rune) {
	sendWhisperIq := handlers.NewInlineQuery(sendwhisperFilter, sendWhisperResponse)
	chosenWhisperIq := handlers.NewChosenInlineResult(chosenWhisperFilter, chosenWhisperResponse)

	d.AddHandler(chosenWhisperIq)
	d.AddHandler(sendWhisperIq)
}
