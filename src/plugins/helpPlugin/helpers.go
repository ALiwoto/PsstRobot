package helpPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func getMainMenuHelpButtons() *gotgbot.InlineKeyboardMarkup {
	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "📖 Privacy Policy",
					CallbackData: "privacy",
				},
				{
					Text:         "📝 Whisper history",
					CallbackData: userWhisperHistoryData,
				},
			},
			{
				{
					Text:              "🧾 Try inline",
					SwitchInlineQuery: new(string),
				},
				{
					Text: "🚑 Support group",
					Url:  "https://t.me/KaizokuBots",
				},
			},
			{
				{
					Text: "☠️kaizoku",
					Url:  "https://t.me/Kaizoku/158",
				},
			},
		},
	}
}

func getUserWhisperHistoryButtons() *gotgbot.InlineKeyboardMarkup {
	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Clear whisper history",
					CallbackData: clearUserHistoryData,
				},
				{
					Text:         "Disable whisper history",
					CallbackData: disableUserHistoryData,
				},
			},
			{
				{
					Text:         "🔙 Back",
					CallbackData: helpCommand,
				},
			},
		},
	}
}
