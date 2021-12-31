package whisperPlugin

import (
	"strconv"
	"strings"

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
		Caption: strings.ReplaceAll(message.Caption, whisperDatabase.CaptionSep, " "),
	}
	switch m.MediaType {
	case whisperDatabase.WhisperTypePlainText:
		e.Caption = strings.ReplaceAll(message.Text, whisperDatabase.CaptionSep, " ")
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

	e.FileId = strings.ReplaceAll(e.FileId, whisperDatabase.MediaGroupSep, " ")

	m.Elements = append(m.Elements, e)
}

func (m *MediaGroupWhisper) getFileIDs() string {
	var result string
	for _, current := range m.Elements {
		result += current.FileId + whisperDatabase.MediaGroupSep
	}
	return result
}

func (m *MediaGroupWhisper) getCaptions() string {
	var result string
	for _, current := range m.Elements {
		result += current.Caption + whisperDatabase.CaptionSep
	}
	return result
}

func (m *MediaGroupWhisper) toWhisper(a *AdvancedWhisper) *whisperDatabase.Whisper {
	w := &whisperDatabase.Whisper{
		Sender:            a.OwnerId,
		Text:              m.getCaptions(),
		Type:              m.MediaType,
		FileId:            m.getFileIDs(),
		Recipient:         a.TargetId,
		RecipientUsername: a.TargetUsername,
	}

	w.GenerateUniqueID()

	return w
}

//---------------------------------------------------------

func (a *AdvancedWhisper) ToWhisper() *whisperDatabase.Whisper {
	if a.MediaGroup != nil {
		return a.MediaGroup.toWhisper(a)
	}
	/*
		w := &Whisper{
			Sender:          result.From.Id,
			Text:            result.Query,
			InlineMessageId: result.InlineMessageId,
			Type:            WhisperTypePlainText,
		}
		w.ParseRecipient(result)
		w.GenerateUniqueID()
		AddWhisper(w)
	*/
	w := &whisperDatabase.Whisper{
		Sender:            a.OwnerId,
		Text:              a.Text,
		Type:              a.MediaType,
		FileId:            a.FileId,
		Recipient:         a.TargetId,
		RecipientUsername: a.TargetUsername,
	}

	w.GenerateUniqueID()

	// we shouldn't use AddWhisper function here.
	//whisperDatabase.AddWhisper(w)

	return w
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
