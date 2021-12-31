package whisperPlugin

import (
	"strings"
	"time"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
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

func cancelWhisperCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return cq.Data == CancelWhisperData
}

func cancelWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	data := usersDatabase.GetUserData(ctx.EffectiveUser.Id)
	if data == nil || data.IsIdle() {
		_, _ = message.Reply(
			bot,
			"No active command to cancel. I wasn't doing anything anyway. Zzzzz...",
			nil,
		)
		return ext.EndGroups
	}

	data.SetToIdle()
	usersDatabase.UpdateUserData(data)

	md := mdparser.GetNormal("The latest operation has been cancelled. ")
	md.AppendNormalThis("Anything else I can do for you?")
	md.AppendNormalThis("\n\nSend /help for a list of possible commands.")

	_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: core.MarkdownV2,
	})

	return ext.EndGroups
}

func generatorListenerFilter(msg *gotgbot.Message) bool {
	return isPrivate(msg) && usersDatabase.IsUserCreating(msg.From)
}

func generatorListenerHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage
	text := strongStringGo.SplitN(message.Text, 2, " ", "\n")
	invalidInput := func() {
		md := mdparser.GetNormal("Invalid username or user-id provided.")
		md.AppendNormalThis("\nTry entering a valid and correct username or user-id.")
		md.AppendNormalThis("\nYou can also enter \"everyone\" if you want everyone")
		md.AppendNormalThis("to be able to see the whisper.")
		_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
		})
	}
	addToMap := func(advanced *AdvancedWhisper) {
		advancedWhisperMutex.Lock()
		advancedWhisperMap[advanced.OwnerId] = advanced
		advancedWhisperMutex.Unlock()
		md := mdparser.GetNormal("Done! This whisper is going to be sent to ")
		md = md.AppendThis(advanced.GetTargetAsMd())
		md.AppendNormalThis(". \nPlease send the content to be whispered. ")
		md.AppendNormalThis("It can be text, photo, or any other media.")
		_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:             core.MarkdownV2,
			DisableWebPagePreview: true,
			ReplyMarkup:           titleChosenMarkup,
		})
	}

	if text[0][1:] == CancelWhisperData {
		return cancelWhisperResponse(bot, ctx)
	}

	advancedWhisperMutex.Lock()
	advanced := advancedWhisperMap[user.Id]
	advancedWhisperMutex.Unlock()

	if advanced == nil {
		if len(text) < 2 {
			invalidInput()
			return ext.EndGroups
		}
		// user still has to enter target's username or user-id.
		advanced = &AdvancedWhisper{
			bot:     bot,
			ctx:     ctx,
			OwnerId: user.Id,
		}

		if isEveryone(text[1]) {
			addToMap(advanced)
			return ext.EndGroups
		}

		username := utils.ExtractUsername(text[1])
		if username != "" {
			advanced.TargetUsername = username
		} else {
			advanced.TargetId = utils.ExtractUserIdFromMessage(message, text[1])
			if advanced.TargetId <= 0 {
				invalidInput()
				return ext.EndGroups
			}
		}

		addToMap(advanced)
		return ext.EndGroups
	}

	if ctx.Message.MediaGroupId != "" {
		// user wants to send a media group whisper.
		return mediaGroupListenerHandler(advanced, ctx)
	}

	advanced.MediaType = extractMediaType(advanced.ctx.Message)
	if advanced.MediaType == whisperDatabase.WhisperTypePlainText {
		advanced.Text = message.Text
	} else {
		advanced.Text = message.Caption
	}
	advanced.FileId = extractFileId(message)

	sendAdvancedWhisperResponse(advanced)

	return ext.ContinueGroups
}

func mediaGroupListenerHandler(advanced *AdvancedWhisper, ctx *ext.Context) error {
	user := advanced.ctx.EffectiveUser
	if advanced.MediaGroup == nil {
		advanced.MediaGroup = &MediaGroupWhisper{
			OwnerId:   user.Id,
			MediaType: extractMediaType(advanced.ctx.Message),
		}
		go sendAdvancedWhisperResponse(advanced)
	}

	advanced.MediaGroup.AddElement(advanced.ctx.Message)

	return ext.ContinueGroups
}

func sendAdvancedWhisperResponse(w *AdvancedWhisper) {
	// sleep for at least 2 seconds, so bot can receive all media group messages
	// from telegram.
	time.Sleep(time.Second * 2)

	advancedWhisperMutex.Lock()
	delete(advancedWhisperMap, w.OwnerId)
	advancedWhisperMutex.Unlock()

	message := w.ctx.Message
	whisper := w.ToWhisper()

	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(
		markup.InlineKeyboard[0], whisper.GetInlineShareButton(),
	)

	whisperDatabase.AddWhisper(whisper)

	md := mdparser.GetNormal("Done! The whisper is ready to be sent.")
	md.AppendNormalThis(". \nClick the button below to share it.")

	_, _ = message.Reply(w.bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		AllowSendingWithoutReply: true,
		DisableWebPagePreview:    true,
		ReplyMarkup:              markup,
	})
}

//---------------------------------------------------------

func createHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage

	if usersDatabase.IsUserBanned(user) {
		return ext.EndGroups
	}

	usersDatabase.ChangeUserStatus(user, usersDatabase.UserStatusCreating)

	md := mdparser.GetNormal("Preparing an advanced whisper. To cancel, type /cancel.")
	md.AppendNormalThis("\n\nSend me the target user's username or id.")
	md.AppendNormalThis("\nYou can also send \"everyone\" to send a whisper to everyone")
	md.AppendNormalThis(" in a chat.")

	_, _ = message.Reply(bot, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                core.MarkdownV2,
		AllowSendingWithoutReply: true,
		DisableWebPagePreview:    true,
	})

	return ext.EndGroups
}

//---------------------------------------------------------
