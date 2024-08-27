package whisperPlugin

import (
	wv "github.com/ALiwoto/PsstRobot/src/core/wotoValues"
	"github.com/ALiwoto/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	generatingInputMessageContent = &gotgbot.InputTextMessageContent{
		MessageText:        "Generating whisper message...",
		LinkPreviewOptions: wv.DisabledWebPagePreview,
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
