package whisperPlugin

import (
	"github.com/AnimeKaizoku/PsstRobot/src/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type MediaGroupWhisper struct {
	OwnerId   int64
	Elements  []*MediaGroupElement
	MediaType whisperDatabase.WhisperType
}

type AdvancedWhisper struct {
	OwnerId        int64
	TargetId       int64
	TargetUsername string
	Text           string
	FileId         string
	bot            *gotgbot.Bot
	ctx            *ext.Context
	MediaGroup     *MediaGroupWhisper
	MediaType      whisperDatabase.WhisperType
}

type MediaGroupElement struct {
	Caption string
	FileId  string
}
