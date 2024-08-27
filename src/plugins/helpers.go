package plugins

import (
	"errors"
	"fmt"

	"github.com/ALiwoto/PsstRobot/src/core/logging"
	"github.com/ALiwoto/PsstRobot/src/core/wotoConfig"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartTelegramBot() error {
	token := wotoConfig.GetBotToken()
	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: 6 * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		MaxRoutines: -1,
	})

	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true, // we don't want the bot to spam people's pm after being down for a while
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			AllowedUpdates: []string{
				gotgbot.UpdateTypeMessage,
				gotgbot.UpdateTypeEditedMessage,
				gotgbot.UpdateTypeChannelPost,
				// gotgbot.UpdateTypeEditedChannelPost,
				gotgbot.UpdateTypeInlineQuery,
				gotgbot.UpdateTypeChosenInlineResult,
				gotgbot.UpdateTypeCallbackQuery,
				gotgbot.UpdateTypeShippingQuery,
				gotgbot.UpdateTypePreCheckoutQuery,
				// gotgbot.UpdateTypePoll,
				// gotgbot.UpdateTypePollAnswer,
				// gotgbot.UpdateTypeMyChatMember,
				// gotgbot.UpdateTypeChatMember,
				// gotgbot.UpdateTypeChatJoinRequest,
			},
			Timeout: 30,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: 6 * gotgbot.DefaultTimeout,
			},
		},
	})
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	LoadAllHandlers(dispatcher, wotoConfig.GetCmdPrefixes())

	updater.Idle()
	return nil
}
