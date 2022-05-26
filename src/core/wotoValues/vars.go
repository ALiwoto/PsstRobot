package wotoValues

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var Core *WotoCore = &WotoCore{}

// handlers that need to be shared globally between all plugins
var (
	CreateWhisperHandler handlers.Response
)
