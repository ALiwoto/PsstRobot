package whisperPlugin

import (
	"sync"

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
					Text:         "🔍 Cancel",
					CallbackData: CancelWhisperData,
				},
			},
		},
	}
)

var (
	advancedWhisperMap   = make(map[int64]*AdvancedWhisper)
	advancedWhisperMutex = &sync.Mutex{}
)
