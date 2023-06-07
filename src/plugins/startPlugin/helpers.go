package startPlugin

import "github.com/PaulSonOfLars/gotgbot/v2"

func getNormalStartButtons() *gotgbot.InlineKeyboardMarkup {
	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:              "📝 Use me",
					SwitchInlineQuery: new(string),
				},
				{
					Text:         "📗 Help",
					CallbackData: helpCommand,
				},
			},
		},
	}
}
