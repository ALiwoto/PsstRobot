package tests

import (
	"log"
	"strings"
	"testing"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
)

func TestSplitting(t *testing.T) {
	const c = whisperDatabase.CaptionSep
	const myStr = "hello" + c + "" + c + "world"
	myStrs := strings.Split(myStr, c)
	log.Println(myStrs)
}
