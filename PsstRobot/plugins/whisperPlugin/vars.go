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
)

var (
	MediaGroupWhisperMap = make(map[int64]*MediaGroupWhisper)
	MediaGroupMutex      = &sync.Mutex{}
)
