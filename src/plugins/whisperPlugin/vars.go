package whisperPlugin

import (
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	generatingInputMessageContent = &gotgbot.InputTextMessageContent{
		MessageText:           "Generating whisper message...",
		DisableWebPagePreview: true,
	}
	titleChosenMarkup = &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "💢 Cancel",
					CallbackData: CancelWhisperData,
				},
			},
		},
	}
)

var (
	advancedWhisperMap = ssg.NewSafeMap[int64, AdvancedWhisper]()
)
