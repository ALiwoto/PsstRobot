package whisperPlugin

import (
	"strconv"

	"github.com/ALiwoto/mdparser/mdparser"
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

//---------------------------------------------------------

func (a *AdvancedWhisper) ToWhisper() *whisperDatabase.Whisper {
	return nil
}

func (a *AdvancedWhisper) IsForEveryone() bool {
	return a.TargetId == 0 && a.TargetUsername == ""
}

func (a *AdvancedWhisper) GetTargetAsMd() mdparser.WMarkDown {
	if a.IsForEveryone() {
		return mdparser.GetBold("everyone")
	}

	if a.TargetUsername != "" {
		return mdparser.GetBold(a.TargetUsername)
	}

	if a.TargetId > 0 {
		if a.bot == nil {
			return mdparser.GetMono(strconv.FormatInt(a.TargetId, 10))
		}
		chat, err := a.bot.GetChat(a.TargetId)
		if err != nil || chat == nil {
			return mdparser.GetMono(strconv.FormatInt(a.TargetId, 10))
		}

		return mdparser.GetUserMention(chat.FirstName, chat.Id)
	}

	/* impossible to reach */
	return mdparser.GetBold("unknown")
}

//---------------------------------------------------------
//---------------------------------------------------------
