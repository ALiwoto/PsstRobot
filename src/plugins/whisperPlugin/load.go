package whisperPlugin

import (
	wv "github.com/ALiwoto/PsstRobot/src/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// LoadHandlers helper function will load all handlers for the current plugin.
func LoadHandlers(d *ext.Dispatcher, t []rune) {
	wv.CreateWhisperHandler = createHandler

	cancelCmd := handlers.NewCommand(CancelWhisperData, cancelWhisperResponse)
	sendWhisperIq := handlers.NewInlineQuery(sendWhisperFilter, sendWhisperResponse)
	chosenWhisperIq := handlers.NewChosenInlineResult(chosenWhisperFilter, chosenWhisperResponse)
	showWhisperCb := handlers.NewCallback(showWhisperCallBackQuery, showWhisperResponse)
	cancelWhisperCb := handlers.NewCallback(cancelWhisperCallBackQuery, cancelWhisperResponse)
	whisperGeneratorListener := handlers.NewMessage(generatorListenerFilter, generatorListenerHandler)
	createCmd := handlers.NewCommand(createCommand, createHandler)

	cancelCmd.Triggers = t
	createCmd.Triggers = t

	d.AddHandler(cancelCmd)
	d.AddHandler(chosenWhisperIq)
	d.AddHandler(sendWhisperIq)
	d.AddHandler(showWhisperCb)
	d.AddHandler(cancelWhisperCb)
	d.AddHandler(whisperGeneratorListener)
	d.AddHandler(createCmd)
}
