package whisperPlugin

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/logging"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

/*
var results []gotgbot.InlineQueryResult

	uname, err := database.GetLastFMUserFromDB(user.Id)
	// fmt.Println(uname)
	if err != nil || uname.LastFmUsername == "" {
		m = mdparser.GetBold(user.FirstName).AppendNormal(" ").AppendItalic("haven't registered themselves on this bot yet.").AppendNormal("\n")
		m = m.AppendBold("Please use ").AppendMono("/setusername")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s's needs to register themselves", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	d, err := lastfm.GetRecentTracksByUsername(uname.LastFmUsername, 2)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	if d.Error != 0 {
		m = mdparser.GetItalic(fmt.Sprintf("Error: %s", d.Message))
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("Sorry %s!, I encountered an error :/", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	if d.Recenttracks == nil || len(d.Recenttracks.Track) == 0 {
		m = mdparser.GetItalic("No tracks were being played.")
		results = append(results, gotgbot.InlineQueryResultArticle{Id: ctx.InlineQuery.Id, Title: fmt.Sprintf("%s haven't played anything yet.", user.FirstName),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2"}})
	}

	lfmUser, err := lastfm.GetLastFMUser(uname.LastFmUsername)
	if err != nil {
		logging.SUGARED.Warn(err.Error())
		return err
	}

	for e, i := range d.Recenttracks.Track {
		var s string
		if d.Recenttracks.Track[e].Attr != nil {
			s = "is now"
		} else {
			s = "was"
		}
		setting, err := database.GetBotUserByID(user.Id)
		if err != nil {
			logging.SUGARED.Warn(err.Error())
			return err
		}

		var m mdparser.WMarkDown
		switch setting.ShowProfile {
		case true:
			m = mdparser.GetHyperLink(user.FirstName, lfmUser.User.URL).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		case false:
			m = mdparser.GetBold(user.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		default:
			m = mdparser.GetBold(user.FirstName).AppendNormal(fmt.Sprintf(" %s listening to", s)).AppendNormal("\n")
		} //mdparser.GetNormal(fmt.Sprintf("%s %s listening to", user.FirstName, s)).AppendNormal("\n")
		m = m.AppendItalic(i.Artist.Name).AppendNormal(" - ").AppendBold(i.Name).AppendNormal("\n")
		m = m.AppendItalic(fmt.Sprintf("%s total plays", lfmUser.User.Playcount))
		if i.Loved == "1" {
			m = m.AppendNormal(", ").AppendItalic("Loved ♥")
		}

		/*if strings.Contains(query.Query, "lyrics") {
			m = m.AppendNormal("\n\n").AppendBold("Lyrics").AppendNormal("\n")
			l, err := LyricsClient.Search(i.Artist.Name, i.Name)
			if err != nil {
				m = m.AppendItalic(fmt.Sprintf("Error: %s", err.Error()))
			}
			m = m.AppendItalic(l)
		} * /
		results = append(results, gotgbot.InlineQueryResultArticle{Id: uuid.New().String(), Title: fmt.Sprintf("%s - %s", i.Artist.Name, i.Name),
			InputMessageContent: gotgbot.InputTextMessageContent{MessageText: m.ToString(), ParseMode: "markdownv2", DisableWebPagePreview: true}})

		if e > 12 {
			break
		}
	}

	if err != nil {
		_, err := query.Answer(b,
			results,
			&gotgbot.AnswerInlineQueryOpts{})
		if err != nil {
			return err
		}
	}
	_, err = query.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{IsPersonal: true})
	if err != nil {
		logging.SUGARED.Error(err.Error())
		return err
	}
	return nil



*/

func sendWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.InlineQuery

	var results []gotgbot.InlineQueryResult
	markup := &gotgbot.InlineKeyboardMarkup{}
	markup.InlineKeyboard = make([][]gotgbot.InlineKeyboardButton, 1)
	markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "🔐 show message",
		CallbackData: "show_whisper",
	})

	results = append(results, &gotgbot.InlineQueryResultArticle{
		Title: "Send a whisper message",
		InputMessageContent: &gotgbot.InputTextMessageContent{
			MessageText:           "Send a whisper message to someone with their username or id",
			ParseMode:             "markdownv2",
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
	print("here")
	return true
}

func chosenWhisperResponse(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = bot.EditMessageText("edit test", &gotgbot.EditMessageTextOpts{
		InlineMessageId: ctx.ChosenInlineResult.InlineMessageId,
	})

	// don't let another handlers to be executed
	return ext.EndGroups
}
