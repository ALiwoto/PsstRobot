package whisperDatabase

import "sync"

var (
	ModelWhisper = &Whisper{}

	// whispersMap is a map with unique id as key and whisper as value.
	whispersMap   = make(map[string]*Whisper)
	whispersMutex = &sync.Mutex{}
)
