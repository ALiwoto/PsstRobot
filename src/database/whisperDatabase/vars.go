package whisperDatabase

import (
	"github.com/ALiwoto/ssg/ssg"
)

// database models
var (
	ModelWhisper = &Whisper{}
)

// caching
var (
	// whispersMap is a map with unique id as key and whisper as value.
	whispersMap = ssg.NewSafeEMap[string, Whisper]()
)
