package whisperDatabase

import (
	"strconv"
	"time"
)

func (w *Whisper) IsExpired(d time.Duration) bool {
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
		return false
	}

	return time.Since(w.CreatedAt) > d
}

func (w *Whisper) GetDBIndex() int {
	return int(strconv.FormatInt(w.Sender, 10)[0] - '0')
}

func (w *Whisper) GenerateUniqueID() {
	if len(w.UniqueId) != 0 {
		return
	}
	chat := strconv.FormatInt(w.ChatId, 10)
	message := strconv.FormatInt(w.MessageId, 10)

	w.UniqueId = chat + "!" + message
}
