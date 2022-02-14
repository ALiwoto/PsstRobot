package tests

import (
	"log"
	"strings"
	"testing"

	"github.com/AnimeKaizoku/PsstRobot/src/database/whisperDatabase"
)

func TestSplitting(t *testing.T) {
	const c = whisperDatabase.CaptionSep
	const myStr = "hello" + c + "" + c + "world"
	myStrs := strings.Split(myStr, c)
	log.Println(myStrs)
}

func TestUserHistorySlices(t *testing.T) {
	history := []string{
		"h1", // remove
		"h2", // remove
		"h3", // remove
		"h4", // keep
		"h5", // keep
		"h6", // keep
	}
	counter := len(history) - 3
	correct := history[counter:]
	removing := history[:counter]
	log.Println(correct, removing)
}
