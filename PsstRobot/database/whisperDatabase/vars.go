package whisperDatabase

import "sync"

// database models
var (
	ModelWhisper = &Whisper{}
)

// caching
var (
	// whispersMap is a map with unique id as key and whisper as value.
	whispersMap   = make(map[string]*Whisper)
	whispersMutex = &sync.Mutex{}
)
