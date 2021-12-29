package whisperPlugin

import (
	"strings"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/logging"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/usersDatabase"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

//---------------------------------------------------------

func showWhisperCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, ShowWhisperData+sepChar)
}

func showWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	user := ctx.EffectiveUser
	myStrs := strings.Split(query.Data, sepChar)
	if len(myStrs) < 2 {
		// cached time shouldn't be long here
		// generating whispers won't take that much long, so cached time
		// should be short enough that user can see the whisper message soon.
		// also, bot will edit the message and change the callback query data
		// anyway, so there is no point to choose it long.
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This whisper is not generated yet...",
			ShowAlert: true,
			CacheTime: 5,
		})
		return ext.EndGroups
	}

	w := whisperDatabase.GetWhisper(myStrs[1])
	if w == nil {
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This whisper is too old...",
			ShowAlert: true,
			CacheTime: 2500,
		})
		return ext.EndGroups
	}

	if w.CanRead(user) {
		if w.ShouldRedirect() {
			// redirect user to bot's pm
			_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
				Url:       w.GetUrl(bot),
				ShowAlert: true,
				CacheTime: 2500,
			})
			return ext.EndGroups
		}

		_, err := query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      w.Text,
			ShowAlert: true,
			CacheTime: 2500,
		})
		if err != nil {
			logging.Error(err)
			return ext.EndGroups
		}

		privacy := userHasPrivacy(w, user)

		if w.ShouldMarkAsRead(user) {
			if !privacy {
				md := mdparser.GetUserMention(user.FirstName, user.Id)
				md.AppendNormalThis(" read the whisper")
				_, _ = bot.EditMessageText(md.ToString(), &gotgbot.EditMessageTextOpts{
					InlineMessageId:       query.InlineMessageId,
					ParseMode:             core.MarkdownV2,
					DisableWebPagePreview: true,
				})
				usersDatabase.SaveInHistory(w.Sender, user)
			}

			whisperDatabase.RemoveWhisper(w)
		}

		return ext.EndGroups
	}

	_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "This psst is not for you!",
		ShowAlert: true,
		CacheTime: 2500,
	})
	return ext.EndGroups
}

//---------------------------------------------------------

func sendwhisperFilter(iq *gotgbot.InlineQuery) bool {
	return true
}

func sendWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery
	user := ctx.EffectiveUser

	if len(query.Query) > MaxQueryLength {
		return answerForLongAdvanced(bot, ctx)
	}

	var results []gotgbot.InlineQueryResult
	var isAnyone bool
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "🔐 show message",
		CallbackData: ShowWhisperData + sepChar,
	})

	var title string
	var description string
	if len(query.Query) > 2 {
		result := utils.ExtractRecipient(query.Query)
		if !result.IsUsernameValid() {
			return answerForHelp(bot, ctx)
		}

		if result.IsForEveryone() {
			isAnyone = true
			title = "📖 A whisper message to anyone!"
			description = "Everyone will be able to open this whisper."
		} else if result.TargetID > 0 {
			var chat *gotgbot.Chat
			var err error
			chat, err = bot.GetChat(result.TargetID)
			if err == nil && chat != nil {
				title = "🔐 A whisper message to " + chat.FirstName
			}
		} else if result.Username != "" {
			title = "🔐 A whisper message to " + result.Username
		}
	} else {
		// seems like an invalid whisper...
		// for the time being we will just show redirect button
		// to them.
		return answerForHelp(bot, ctx)
	}

	if title == "" {
		title = "📍Send a whisper message to someone with their username or id"
	}

	if description == "" {
		description = "Only they can open this whisper."
	}

	results = append(results, &gotgbot.InlineQueryResultArticle{
		Id:                  ctx.InlineQuery.Id,
		Title:               title,
		Description:         description,
		InputMessageContent: generatingInputMessageContent,
		ReplyMarkup:         markup,
	})

	if isAnyone {
		history := usersDatabase.GetUserHistory(user.Id)
		if history != nil && !history.IsEmpty() {
			timeStampNow := utils.ToBase32(time.Now().Unix())
			var resultId string
			for _, current := range history.History {
				// time::user::target
				resultId = timeStampNow
				resultId += wv.ResultIdentifier + utils.ToBase32(current.OwnerId)
				resultId += wv.ResultIdentifier + utils.ToBase10(current.TargetId)
				title = "📩 A whisper message to " + current.TargetName
				description = "Only " + current.TargetName + " can open this whisper."
				results = append(results, &gotgbot.InlineQueryResultArticle{
					Id:                  resultId,
					Title:               title,
					Description:         description,
					InputMessageContent: generatingInputMessageContent,
					ReplyMarkup:         markup,
				})
			}

		}
	}

	_, _ = query.Answer(bot, results, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
	})

	// don't let another handlers to be executed
	return ext.EndGroups
}

//---------------------------------------------------------

func chosenWhisperFilter(cir *gotgbot.ChosenInlineResult) bool {
	return true
}

func chosenWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	result := ctx.ChosenInlineResult
	w := whisperDatabase.CreateNewWhisper(result)
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "🔐 Show message",
		CallbackData: ShowWhisperData + sepChar + w.UniqueId,
	})

	_, _ = bot.EditMessageText(w.ParseAsMd(bot).ToString(), &gotgbot.EditMessageTextOpts{
		ReplyMarkup:     *markup,
		ParseMode:       core.MarkdownV2,
		InlineMessageId: result.InlineMessageId,
	})

	// don't let another handlers to be executed
	return ext.EndGroups
}

//---------------------------------------------------------

func generatorListenerFilter(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "private" && usersDatabase.IsUserCreating(msg.From)
}

func generatorListenerHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.Message.MediaGroupId != "" {
		// user wants to send a media group whisper.
		return mediaGroupListenerHandler(bot, ctx)
	}

	//message := ctx.Message
	return ext.ContinueGroups
}

func mediaGroupListenerHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	MediaGroupMutex.Lock()
	w := MediaGroupWhisperMap[user.Id]
	MediaGroupMutex.Unlock()
	if w == nil {
		w = &MediaGroupWhisper{
			OwnerId:   user.Id,
			MediaType: extractMediaType(ctx.Message),
			ctx:       ctx,
			bot:       bot,
		}
		MediaGroupMutex.Lock()
		MediaGroupWhisperMap[user.Id] = w
		MediaGroupMutex.Unlock()
		go sendMediaGroupResponse(w)
	}

	w.AddElement(ctx.Message)

	return ext.ContinueGroups
}

func sendMediaGroupResponse(w *MediaGroupWhisper) {
	// sleep for at least 2 seconds, so bot can receive all media group messages
	// from telegram.
	time.Sleep(time.Second * 2)

	MediaGroupMutex.Lock()
	delete(MediaGroupWhisperMap, w.OwnerId)
	MediaGroupMutex.Unlock()

	message := w.ctx.Message
	whisper := w.ToWhisper()

	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(
		markup.InlineKeyboard[0], whisper.GetInlineShareButton(),
	)

	whisperDatabase.AddWhisper(whisper)

	_, _ = message.Reply(w.bot, "", &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		AllowSendingWithoutReply: true,
		DisableWebPagePreview:    true,
		ReplyMarkup:              markup,
	})
}

//---------------------------------------------------------
