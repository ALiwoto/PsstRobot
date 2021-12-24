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

func (w *Whisper) ParseAsMd(bot *gotgbot.Bot) mdparser.WMarkDown {
	md := mdparser.GetNormal("A whisper message to ")
	var rec mdparser.WMarkDown
	if w.RecipientUsername != "" {
		rec = mdparser.GetNormal(w.RecipientUsername)
	} else if w.Recipient != 0 {
		chat, _ := bot.GetChat(w.Recipient)
		if chat != nil {
			rec = mdparser.GetUserMention(chat.FirstName, chat.Id)
		} else {
			rec = mdparser.GetMono(strconv.FormatInt(w.Recipient, 10))
		}
	}

	if rec != nil {
		md.AppendThis(rec)
		md.AppendNormalThis(".\nOnly they can read the message.")
	} else {
		md.AppendNormalThis("anyone.")
		md.AppendNormalThis("\nAnyone can read it!")
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

	if u.Id == w.Sender {
		return true
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

func (w *Whisper) ShouldMarkAsRead(u *gotgbot.User) bool {
	return u.Id != w.Sender
}

func (w *Whisper) setText(value string) {
	w.Text = strings.TrimSpace(value)
}

func (w *Whisper) Unpack() (*utils.UnpackInlineMessageResult, error) {
	if w.unpackedResult != nil {
		return w.unpackedResult, nil
	}

	r, err := utils.UnpackInlineMessageId(w.InlineMessageId)
	if err != nil {
		return nil, err
	}
	w.unpackedResult = r

	return r, nil
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
