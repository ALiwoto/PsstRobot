package whisperDatabase

import "time"

type Whisper struct {
	UniqueId          string    `json:"unique_id" gorm:"primaryKey"`
	InlineMessageId   string    `json:"inline_message_id"`
	Sender            int64     `json:"sender"`
	Text              string    `json:"text"`
	Recipient         int64     `json:"recipient"`
	RecipientUsername string    `json:"recipient_username"`
	CreatedAt         time.Time `json:"created_at"`
}
