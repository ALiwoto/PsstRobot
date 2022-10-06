package helpPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// LoadHandlers helper function will load all handlers for the current plugin.
func LoadHandlers(d *ext.Dispatcher, t []rune) {
	helpCmd := handlers.NewCommand(helpCommand, helpHandler)

	helpCmd.Triggers = t

	d.AddHandler(helpCmd)
}
