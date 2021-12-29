package startPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// LoadHandlers helper function will load all handlers for the current plugin.
func LoadHandlers(d *ext.Dispatcher, t []rune) {
	startCmd := handlers.NewCommand(startCommand, startHandler)
	privacyCmd := handlers.NewCommand(privacyCommand, privacyHandler)
	startCmd.Triggers = t
	privacyCmd.Triggers = t

	d.AddHandler(startCmd)
	d.AddHandler(privacyCmd)
}
