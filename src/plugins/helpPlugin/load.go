package helpPlugin

import (
	wv "github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// LoadHandlers helper function will load all handlers for the current plugin.
func LoadHandlers(d *ext.Dispatcher, t []rune) {
	wv.HelpHandler = helpHandler
	helpCmd := handlers.NewCommand(helpCommand, helpHandler)
	userWhisperHistoryCmd := handlers.NewCommand(userWhisperHistoryData, userWhisperHistoryResponse)
	userWhisperHistoryCb := handlers.NewCallback(userWhisperHistoryCallBackQuery, userWhisperHistoryResponse)
	helpCb := handlers.NewCallback(helpCallBackQuery, helpHandler)
	clearUserHistoryCb := handlers.NewCallback(clearUserHistoryCallBackQuery, clearUserHistoryResponse)
	disableUserHistoryCb := handlers.NewCallback(disableUserHistoryCallBackQuery, disableUserHistoryResponse)

	helpCmd.Triggers = t
	userWhisperHistoryCmd.Triggers = t

	d.AddHandler(helpCmd)
	d.AddHandler(helpCb)
	d.AddHandler(userWhisperHistoryCmd)
	d.AddHandler(userWhisperHistoryCb)
	d.AddHandler(clearUserHistoryCb)
	d.AddHandler(disableUserHistoryCb)
}
