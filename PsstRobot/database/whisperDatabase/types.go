package whisperDatabase

import (
	"time"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
)

type WhisperType int

type Whisper struct {
	UniqueId          string                           `json:"unique_id" gorm:"primaryKey"`
	InlineMessageId   string                           `json:"inline_message_id"`
	Sender            int64                            `json:"sender"`
	Text              string                           `json:"text"`
	Recipient         int64                            `json:"recipient"`
	RecipientUsername string                           `json:"recipient_username"`
	FileId            string                           `json:"file_id"`
	Type              WhisperType                      `json:"type"`
	CaptionIndex      int                              `json:"caption_index"`
	CreatedAt         time.Time                        `json:"created_at"`
	unpackedResult    *utils.UnpackInlineMessageResult `json:"-" gorm:"-" sql:"-"`
}
