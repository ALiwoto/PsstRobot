package whisperDatabase

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
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
	return utils.GetDBIndex(w.Sender)
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

func (w *Whisper) GetInlineTitle(bot *gotgbot.Bot) string {
	var name string
	if w.RecipientUsername != "" {
		name = w.RecipientUsername
	} else if w.Recipient != 0 {
		chat, _ := bot.GetChat(w.Recipient)
		if chat != nil {
			name = chat.FirstName
		} else {
			name = strconv.FormatInt(w.Recipient, 10)
		}
	} else {
		name = "everyone"
	}
	return "An advanced whisper message to " + name + "."
}

func (w *Whisper) GetInlineDescription() string {
	return "Only they can read this advanced whisper."
}

func (w *Whisper) ShouldRedirect() bool {
	return w.Type != WhisperTypePlainText || w.IsTooLong()
}

func (w *Whisper) IsMediaGroup() bool {
	return strings.Contains(w.FileId, MediaGroupSep)
}

func (w *Whisper) GetMediaGroup() []gotgbot.InputMedia {
	var result []gotgbot.InputMedia
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	switch w.Type {
	case WhisperTypeAnimation:
		result = w.getAnimationArray()
	case WhisperTypeDocument:
		result = w.getDocumentArray()
	case WhisperTypeAudio:
		result = w.getAudioArray()
	case WhisperTypePhoto:
		result = w.getPhotoArray()
	case WhisperTypeVideo:
		result = w.getVideoArray()
	}

	// drop the last result
	return result[:len(result)-1]
}

func (w *Whisper) getFileIDs() []string {
	return strings.Split(w.FileId, MediaGroupSep)
}

func (w *Whisper) getCaptions() []string {
	return strings.Split(w.Text, CaptionSep)
}

// getAnimationArray returns an array of InputMediaDocument in form of
// []gotgbot.InputMedia
func (w *Whisper) getDocumentArray() []gotgbot.InputMedia {
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	var myArray []gotgbot.InputMedia
	files := w.getFileIDs()
	captions := w.getCaptions()
	for i, current := range files {
		myArray = append(myArray, gotgbot.InputMediaDocument{
			Media:   current,
			Caption: captions[i],
		})
	}

	return myArray
}

// getAnimationArray returns an array of InputMediaAudio in form of
// []gotgbot.InputMedia
func (w *Whisper) getAudioArray() []gotgbot.InputMedia {
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	var myArray []gotgbot.InputMedia
	files := w.getFileIDs()
	captions := w.getCaptions()
	for i, current := range files {
		myArray = append(myArray, gotgbot.InputMediaAudio{
			Media:   current,
			Caption: captions[i],
		})
	}

	return myArray
}

// getAnimationArray returns an array of InputMediaVideo in form of
// []gotgbot.InputMedia
func (w *Whisper) getVideoArray() []gotgbot.InputMedia {
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	var myArray []gotgbot.InputMedia
	files := w.getFileIDs()
	captions := w.getCaptions()
	for i, current := range files {
		myArray = append(myArray, gotgbot.InputMediaVideo{
			Media:   current,
			Caption: captions[i],
		})
	}

	return myArray
}

// getAnimationArray returns an array of InputMediaPhoto in form of
// []gotgbot.InputMedia
func (w *Whisper) getPhotoArray() []gotgbot.InputMedia {
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	var myArray []gotgbot.InputMedia
	files := w.getFileIDs()
	captions := w.getCaptions()
	for i, current := range files {
		myArray = append(myArray, gotgbot.InputMediaPhoto{
			Media:   current,
			Caption: captions[i],
		})
	}

	return myArray
}

// getAnimationArray returns an array of InputMediaAnimation in form of
// []gotgbot.InputMedia
func (w *Whisper) getAnimationArray() []gotgbot.InputMedia {
	// - InputMediaAnimation
	// - InputMediaDocument
	// - InputMediaAudio
	// - InputMediaPhoto
	// - InputMediaVideo
	var myArray []gotgbot.InputMedia
	files := w.getFileIDs()
	captions := w.getCaptions()
	for i, current := range files {
		myArray = append(myArray, gotgbot.InputMediaAnimation{
			Media:   current,
			Caption: captions[i],
		})
	}

	return myArray
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

		return w.canMatchUsername(u.Username)
	} else if u.Id == w.Recipient {
		return true
	}

	return false
}

func (w *Whisper) CanSendInlineAdvanced(user *gotgbot.User) bool {
	return w.InlineMessageId == "" && user.Id == w.Sender
}

func (w *Whisper) IsForEveryone() bool {
	return w.Recipient == 0 && w.RecipientUsername == ""
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

func (w *Whisper) parseRecipientByResultId(myStrs []string, chosen *gotgbot.ChosenInlineResult) {
	// format:
	// time::user::target

	w.Recipient, _ = strconv.ParseInt(myStrs[2], 10, 64)
	w.setText(chosen.Query)
}

func (w *Whisper) ParseRecipient(chosen *gotgbot.ChosenInlineResult) {
	if strings.Contains(chosen.ResultId, wv.ResultIdentifier) {
		myStrs := strings.Split(chosen.ResultId, wv.ResultIdentifier)
		if len(myStrs) == 3 {
			w.parseRecipientByResultId(myStrs, chosen)
			return
		}
	}

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

func (w *Whisper) GetInlineShareButton() gotgbot.InlineKeyboardButton {
	s := wv.AdvancedInlinePrefix + w.UniqueId + wv.AdvancedInlineSuffix
	return gotgbot.InlineKeyboardButton{
		Text:              "📤 share whisper",
		SwitchInlineQuery: &s,
	}
}
