package whisperDatabase

import (
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

// database models
var (
	ModelWhisper = &Whisper{}
)

// caching
var (
	// whispersMap is a map with unique id as key and whisper as value.
	whispersMap = strongStringGo.NewSafeEMap[string, Whisper]()
)
