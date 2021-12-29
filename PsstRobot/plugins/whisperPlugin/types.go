package whisperPlugin

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type MediaGroupWhisper struct {
	OwnerId   int64
	Elements  []*MediaGroupElement
	MediaType whisperDatabase.WhisperType
	bot       *gotgbot.Bot
	ctx       *ext.Context
}

type MediaGroupElement struct {
	Caption string
	FileId  string
}
