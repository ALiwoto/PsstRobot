package plugins

import (
	"errors"
	"fmt"
	"net/http"

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
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: 6 * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return err
	}

	uOptions := &ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			MaxRoutines: -1,
		}),
	}

	updater := ext.NewUpdater(uOptions)
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: false,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 30,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: 6 * gotgbot.DefaultTimeout,
			},
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
		},
	})
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	LoadAllHandlers(updater.Dispatcher, wotoConfig.GetCmdPrefixes())

	updater.Idle()
	return nil
}
