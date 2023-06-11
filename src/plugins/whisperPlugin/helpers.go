package whisperPlugin

import (
	"strings"

	wv "github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	"github.com/AnimeKaizoku/PsstRobot/src/database/usersDatabase"
	"github.com/AnimeKaizoku/PsstRobot/src/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// answerForLongAdvanced is a helper function to tell user that the whisper is too long
// and they have to use advanced whisper instead of normal whispers.
// it will redirect them to bot's pm with parameter equal to `wv.StartDataCreate`.
// `wv.StartDataCreate` is a constant with value of "create"; so basically,
// user's message to bot will exactly be "/start create".
func answerForLongAdvanced(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.InlineQuery.Answer(bot, nil, &gotgbot.AnswerInlineQueryOpts{
		Button: &gotgbot.InlineQueryResultsButton{
			Text:           "Too long! Use an advanced whisper!",
			StartParameter: wv.StartDataCreate,
		},
	})
	return ext.EndGroups
}

func answerForWrongAdvanced(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.InlineQuery.Answer(bot, nil, &gotgbot.AnswerInlineQueryOpts{
		Button: &gotgbot.InlineQueryResultsButton{
			Text:           "Invalid advanced whisper! Please create another one!",
			StartParameter: wv.StartDataCreate,
		},
	})
	return ext.EndGroups
}

// answerForHelp is a helper function and should be called when the query entered by user
// is wrong. it will redirect the user to bot's pm with parameter equal to `wv.HelpDataInline`.
// `wv.HelpDataInline` is a constant with value of "help-inline"; so basically,
// user's message to bot will exactly be "/start help-inline".
func answerForHelp(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.InlineQuery.Answer(bot, nil, &gotgbot.AnswerInlineQueryOpts{
		Button: &gotgbot.InlineQueryResultsButton{
			Text:           "Invalid whisper! Please create another one!",
			StartParameter: wv.HelpDataInline,
		},
	})
	return ext.EndGroups
}

// userHasPrivacy is a helper function to check if user has privacy mode enabled
// or not.
// if users with privacy mode enabled try to read a whisper, the whisper message
// shouldn't be edited by the bot.
// so we don't expect the bot to edit the message to "{mention} read the whisper".
// but still, even if a user has privacy mode enabled and click on the whisper,
// the whisper should be removed from the database.
// as a result, if a new user click on will get a message that the whisper is removed
// from the db (we know it as "This whisper is too old....").
// but for this issue, we have a solution: caching on client-side.
// `cachedTime` parameter while answering a callback query is set to a long time,
// so even if whisper is removed from db and same users try to click on the button,
// it will be there for them.
func userHasPrivacy(w *whisperDatabase.Whisper, user *gotgbot.User) bool {
	return !w.IsForEveryone() && usersDatabase.HasPrivacy(user)
}

func extractMediaType(message *gotgbot.Message) whisperDatabase.WhisperType {
	// WhisperTypePhoto
	// WhisperTypeVideo
	// WhisperTypeAudio
	// WhisperTypeVoice
	// WhisperTypeSticker
	// WhisperTypeDocument
	// WhisperTypeVideoNote
	// WhisperTypeAnimation
	// WhisperTypeDice
	switch {
	case len(message.Text) > 0:
		return whisperDatabase.WhisperTypePlainText
	case len(message.Photo) > 0:
		return whisperDatabase.WhisperTypePhoto
	case message.Video != nil:
		return whisperDatabase.WhisperTypeVideo
	case message.Audio != nil:
		return whisperDatabase.WhisperTypeAudio
	case message.Voice != nil:
		return whisperDatabase.WhisperTypeVoice
	case message.Sticker != nil:
		return whisperDatabase.WhisperTypeSticker
	case message.Document != nil:
		return whisperDatabase.WhisperTypeDocument
	case message.VideoNote != nil:
		return whisperDatabase.WhisperTypeVideoNote
	case message.Animation != nil:
		return whisperDatabase.WhisperTypeAnimation
	case message.Dice != nil:
		return whisperDatabase.WhisperTypeDice
	}

	return whisperDatabase.WhisperTypeUnknown
}

func extractFileId(message *gotgbot.Message) string {
	// WhisperTypePhoto
	// WhisperTypeVideo
	// WhisperTypeAudio
	// WhisperTypeVoice
	// WhisperTypeSticker
	// WhisperTypeDocument
	// WhisperTypeVideoNote
	// WhisperTypeAnimation
	// WhisperTypeDice
	switch {
	case len(message.Text) > 0:
		return ""
	case len(message.Photo) > 0:
		return message.Photo[0].FileId
	case message.Video != nil:
		return message.Video.FileId
	case message.Audio != nil:
		return message.Audio.FileId
	case message.Voice != nil:
		return message.Voice.FileId
	case message.Sticker != nil:
		return message.Sticker.FileId
	case message.Document != nil:
		return message.Document.FileId
	case message.VideoNote != nil:
		return message.VideoNote.FileId
	case message.Animation != nil:
		return message.Animation.FileId
	case message.Dice != nil:
		return message.Dice.Emoji
	}

	return ""
}

func isPrivate(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "private"
}

func isEveryone(text string) bool {
	return strings.EqualFold(text, "everyone")
}
