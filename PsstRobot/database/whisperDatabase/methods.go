package whisperDatabase

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (w *Whisper) IsExpired(d time.Duration) bool {
	if w.CreatedAt.IsZero() {
		// prevent from possible bugs of removing whisper
		// before the actual expiry date...
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

	part1 := strconv.FormatInt(time.Now().Unix(), 32)
	part2 := strconv.FormatInt(w.Sender, 32)

	w.UniqueId = part1 + "=" + part2
}

func (w *Whisper) ParseAsMd() mdparser.WMarkDown {
	md := mdparser.GetNormal("A whisper message to ")
	var rec string
	if w.RecipientUsername != "" {
		rec = w.RecipientUsername
	} else {
		rec = strconv.FormatInt(w.Recipient, 10)
	}

	if rec != "0" {
		md.AppendNormalThis(rec)
		md.AppendNormalThis(".\nOnly they can read the message.")
	} else {
		md.AppendNormalThis("anyone.")
		md.AppendNormalThis(".\nAnyone can read it!")
	}

	return md
}

func (w *Whisper) canMatchUsername(username string) bool {
	return strings.EqualFold(w.RecipientUsername, "@"+username) ||
		strings.EqualFold(w.RecipientUsername, username)
}

func (w *Whisper) GetUrl(b *gotgbot.Bot) string {
	return "http://t.me/" + b.Username + "?start=" + url.QueryEscape(w.UniqueId)
}

func (w *Whisper) CanRead(u *gotgbot.User) bool {
	if u == nil {
		return false
	}

	// everyone?
	if w.Recipient == 0 {
		if w.RecipientUsername == "" {
			return true
		}

		if w.canMatchUsername(u.Username) {
			return true
		}
	} else if u.Id == w.Recipient {
		return true
	}

	return false
}

func (w *Whisper) setText(value string) {
	w.Text = strings.TrimSpace(value)
}

func (w *Whisper) ParseRecipient() {
	r := utils.ExtractRecipient(w.Text)
	if r == nil {
		return
	}

	w.Recipient = r.TargetID
	w.RecipientUsername = r.Username
	w.setText(r.Text)
}

func (w *Whisper) IsTooLong() bool {
	return len(w.Text) > MaxTextLength
}
