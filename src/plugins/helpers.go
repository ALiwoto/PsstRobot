package plugins

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AnimeKaizoku/PsstRobot/src/core/logging"
	"github.com/AnimeKaizoku/PsstRobot/src/core/wotoConfig"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartTelegramBot() error {
	token := wotoConfig.GetBotToken()
	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  2 * gotgbot.DefaultGetTimeout,
		PostTimeout: 2 * gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		return err
	}

	uOptions := &ext.UpdaterOpts{
		DispatcherOpts: ext.DispatcherOpts{
			MaxRoutines: -1,
		},
	}

	utmp := ext.NewUpdater(uOptions)
	updater := &utmp
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: false,
	})
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	LoadAllHandlers(updater.Dispatcher, wotoConfig.GetCmdPrefixes())

	updater.Idle()
	return nil
}
