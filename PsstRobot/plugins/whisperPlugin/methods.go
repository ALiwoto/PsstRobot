package whisperPlugin

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (m *MediaGroupWhisper) AddElement(message *gotgbot.Message) {
	// WhisperTypePhoto
	// WhisperTypeVideo
	// WhisperTypeAudio
	// WhisperTypeVoice
	// WhisperTypeSticker
	// WhisperTypeDocument
	// WhisperTypeVideoNote
	// WhisperTypeAnimation
	// WhisperTypeDice
	e := &MediaGroupElement{
		Caption: message.Caption,
	}
	switch m.MediaType {
	case whisperDatabase.WhisperTypePlainText:
		e.Caption = message.Text
	case whisperDatabase.WhisperTypePhoto:
		e.FileId = message.Photo[len(message.Photo)-1].FileId
	case whisperDatabase.WhisperTypeVideo:
		e.FileId = message.Video.FileId
	case whisperDatabase.WhisperTypeAudio:
		e.FileId = message.Audio.FileId
	case whisperDatabase.WhisperTypeVoice:
		e.FileId = message.Voice.FileId
	case whisperDatabase.WhisperTypeSticker:
		e.FileId = message.Sticker.FileId
	case whisperDatabase.WhisperTypeDocument:
		e.FileId = message.Document.FileId
	case whisperDatabase.WhisperTypeVideoNote:
		e.FileId = message.VideoNote.FileId
	case whisperDatabase.WhisperTypeAnimation:
		e.FileId = message.Animation.FileId
	case whisperDatabase.WhisperTypeDice:
		e.Caption = message.Dice.Emoji
	}

	m.Elements = append(m.Elements, e)
}

func (m *MediaGroupWhisper) ToWhisper() *whisperDatabase.Whisper {
	return nil
}
