package main

import (
	"log"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/logging"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoConfig"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/plugins"
)

func main() {
	_, err := wotoConfig.LoadConfig()
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}

	f := logging.LoadLogger()
	if f != nil {
		defer f()
	}

	err = plugins.StartTelegramBot()
	if err != nil {
		logging.Fatal("Failed to start the bot bot: ", err)
	}
}
