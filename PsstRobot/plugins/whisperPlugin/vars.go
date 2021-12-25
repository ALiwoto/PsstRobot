package whisperPlugin

import "github.com/PaulSonOfLars/gotgbot/v2"

var (
	generatingInputMessageContent = &gotgbot.InputTextMessageContent{
		MessageText:           "Generating whisper message...",
		DisableWebPagePreview: true,
	}
)
