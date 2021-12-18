package whisperDatabase

import (
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
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

	md.AppendNormalThis(rec)
	md.AppendNormalThis(".\nOnly he/she can read the message.")

	return md
}

func (w *Whisper) canMatchUsername(username string) bool {
	return strings.EqualFold(w.RecipientUsername, "@"+username) ||
		strings.EqualFold(w.RecipientUsername, username)
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

func (w *Whisper) ParseRecipient() {
	if len(w.Text) == 0 {
		return
	}

	// supported formats:
	// @username message
	// message @username
	// ID message
	// message ID
	myStrs := strings.Split(w.Text, " ")
	if myStrs[0][0] == '@' {
		w.RecipientUsername = myStrs[0]
		w.Text = strings.Join(myStrs[1:], " ")
		return
	}

	last := len(myStrs) - 1
	if myStrs[last][0] == '@' {
		w.RecipientUsername = myStrs[last]
		w.Text = strings.Join(myStrs[:last], " ")
		return
	}

	id, err := strconv.ParseInt(myStrs[0], 10, 64)
	if err == nil {
		w.Recipient = id
		w.Text = strings.Join(myStrs[1:], " ")
		return
	}

	id, err = strconv.ParseInt(myStrs[last], 10, 64)
	if err == nil {
		w.Recipient = id
		w.Text = strings.Join(myStrs[1:], " ")
		return
	}

}
