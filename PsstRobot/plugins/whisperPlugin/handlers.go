package whisperPlugin

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func sendWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	// don't let another handlers to be executed
	return ext.EndGroups
}

//func(cir *gotgbot.ChosenInlineResult) bool

func sendwhisperFilter(iq *gotgbot.InlineQuery) bool {
	return false
}

func chosenWhisperFilter(cir *gotgbot.ChosenInlineResult) bool {
	print("here")
	return false
}

func chosenWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	// don't let another handlers to be executed
	return ext.EndGroups
}
