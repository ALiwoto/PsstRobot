package helpPlugin

import (
	"github.com/AnimeKaizoku/PsstRobot/src/database/usersDatabase"
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

func disableOrEnable(value bool) string {
	if value {
		return "Enable"
	}

	return "Disable"
}

func getUserWhisperHistoryButtons(userId int64) *gotgbot.InlineKeyboardMarkup {
	historyDisabled := usersDatabase.IsHistoryDisabled(userId)

	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Clear whisper history",
					CallbackData: clearUserHistoryData,
				},
				{
					Text:         disableOrEnable(!historyDisabled) + " whisper history",
					CallbackData: toggleUserHistoryData,
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
