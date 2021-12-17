package whisperDatabase

import "time"

type Whisper struct {
	UniqueId  string    `json:"unique_id" gorm:"primaryKey"`
	Sender    int64     `json:"sender"`
	ChatId    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	MessageId int64     `json:"message_id"`
	Recipient int64     `json:"recipient"`
	CreatedAt time.Time `json:"created_at"`
}
