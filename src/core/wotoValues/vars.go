package wotoValues

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var Core *WotoCore = &WotoCore{}

// handlers that need to be shared globally between all plugins
var (
	// CreateWhisperHandler is the handler responsible for creating a whisper,
	// set in `path://src/plugins/whisperPlugin/helpers.go`.
	CreateWhisperHandler handlers.Response
	// HelpHandler is the handler responsible for sending the help message
	// to the user.
	// set in `path://src/plugins/helpPlugin/helpers.go`.
	HelpHandler handlers.Response
)

// common variables that can be used everywhere (hopefully)
var (
	DisabledWebPagePreview = &gotgbot.LinkPreviewOptions{
		IsDisabled: true,
	}
)
