package plugins

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/plugins/whisperPlugin"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func LoadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	whisperPlugin.LoadHandlers(d, triggers)
}
