package helpPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func getMainMenuHelpButtons() *gotgbot.InlineKeyboardMarkup {
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 3)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "📖 Privacy Policy",
		CallbackData: "privacy",
	})
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "📝 Clear History",
		CallbackData: "clear_history",
	})
	markup.InlineKeyboard[1] = append(markup.InlineKeyboard[1], gotgbot.InlineKeyboardButton{
		Text:              "🧾 Try inline",
		SwitchInlineQuery: new(string),
	})
	markup.InlineKeyboard[1] = append(markup.InlineKeyboard[1], gotgbot.InlineKeyboardButton{
		Text: "🚑 Support group",
		Url:  "https://t.me/KaizokuBots",
	})
	markup.InlineKeyboard[1] = append(markup.InlineKeyboard[2], gotgbot.InlineKeyboardButton{
		Text: "Made with ❤️ by ☠️kaizoku",
		Url:  "https://t.me/Kaizoku/158",
	})

	return markup
}

// LoadHandlers helper function will load all handlers for the current plugin.
func LoadHandlers(d *ext.Dispatcher, t []rune) {
	helpCmd := handlers.NewCommand(helpCommand, helpHandler)

	helpCmd.Triggers = t

	d.AddHandler(helpCmd)
}
