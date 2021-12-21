package whisperPlugin

import (
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/logging"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func showWhisperCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, ShowWhisperData+sepChar)
}

func showWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	user := ctx.EffectiveUser
	myStrs := strings.Split(query.Data, sepChar)
	if len(myStrs) < 2 {
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
			CacheTime: 5,
		})
		return ext.EndGroups
	}

	if w.Sender == user.Id {
		if w.IsTooLong() {
			// redirect user to bot's pm
			_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
				Url:       w.GetUrl(bot),
				ShowAlert: true,
				CacheTime: 5,
			})
			return ext.EndGroups
		}
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      w.Text,
			ShowAlert: true,
			CacheTime: 5,
		})
		return ext.EndGroups
	}

	if w.CanRead(user) {
		if w.IsTooLong() {
			// redirect user to bot's pm
			_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
				Url:       w.GetUrl(bot),
				ShowAlert: true,
				CacheTime: 5,
			})
			return ext.EndGroups
		}

		go whisperDatabase.RemoveWhisper(w)
		_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text:      w.Text,
			ShowAlert: true,
			CacheTime: 5,
		})
		md := mdparser.GetUserMention(user.FirstName, user.Id)
		md.AppendNormalThis(" read the whisper")
		_, _ = bot.EditMessageText(md.ToString(), &gotgbot.EditMessageTextOpts{
			InlineMessageId:       query.InlineMessageId,
			ParseMode:             "markdownv2",
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}

	_, _ = query.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "This psst is not for you!",
		ShowAlert: true,
		CacheTime: 5,
	})
	return ext.EndGroups
}

func sendWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery

	var results []gotgbot.InlineQueryResult
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "🔐 show message",
		CallbackData: ShowWhisperData + sepChar,
	})

	var title string
	var description string
	if len(query.Query) > 2 && strings.Contains(strings.TrimSpace(query.Query), " ") {
		result := utils.ExtractRecipient(query.Query)
		if result.IsForEveryone() {
			title = "📖 A whisper message to anyone!"
			description = "Everyone will be able to open this whisper."
		} else if result.TargetID > 0 {
			var chat *gotgbot.Chat
			var err error
			chat, err = bot.GetChat(result.TargetID)
			if err == nil && chat != nil {
				title = "🔐 A whisper message to " + chat.Title
			}
		} else if result.Username != "" {
			title = "🔐 A whisper message to " + result.Username
		}

	}

	if title == "" {
		title = "📍Send a whisper message to someone with their username or id"
	}

	if description == "" {
		description = "Only they can open this whisper."
	}

	results = append(results, &gotgbot.InlineQueryResultArticle{
		Id:          ctx.InlineQuery.Id,
		Title:       "Send a whisper message to someone with their username or id",
		Description: "Only they can open this whisper.",
		InputMessageContent: &gotgbot.InputTextMessageContent{
			MessageText:           "Generating whisper message...",
			DisableWebPagePreview: true,
		},
		ReplyMarkup: markup,
	})

	_, err := query.Answer(bot, results, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
	})
	if err != nil {
		logging.Error(err)
	}

	// don't let another handlers to be executed
	return ext.EndGroups
}

//func(cir *gotgbot.ChosenInlineResult) bool

func sendwhisperFilter(iq *gotgbot.InlineQuery) bool {
	return true
}

func chosenWhisperFilter(cir *gotgbot.ChosenInlineResult) bool {
	return true
}

func chosenWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	w := whisperDatabase.CreateNewWhisper(ctx.ChosenInlineResult)
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "🔐 Show message",
		CallbackData: ShowWhisperData + sepChar + w.UniqueId,
	})

	_, _ = bot.EditMessageText(w.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ReplyMarkup:     *markup,
		ParseMode:       "markdownv2",
		InlineMessageId: ctx.ChosenInlineResult.InlineMessageId,
	})

	// don't let another handlers to be executed
	return ext.EndGroups
}
