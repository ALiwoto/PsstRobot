package helpPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
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
	markup.InlineKeyboard[2] = append(markup.InlineKeyboard[2], gotgbot.InlineKeyboardButton{
		Text: "☠️kaizoku",
		Url:  "https://t.me/Kaizoku/158",
	})

	return markup
}
