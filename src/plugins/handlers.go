package plugins

import (
	"github.com/AnimeKaizoku/PsstRobot/src/plugins/helpPlugin"
	"github.com/AnimeKaizoku/PsstRobot/src/plugins/startPlugin"
	"github.com/AnimeKaizoku/PsstRobot/src/plugins/whisperPlugin"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func LoadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	whisperPlugin.LoadHandlers(d, triggers)
	startPlugin.LoadHandlers(d, triggers)
	helpPlugin.LoadHandlers(d, triggers)
}
