package whisperDatabase

import (
	"sync"
	"time"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func LoadAllWhispers() {
	allSessions := wv.Core.GetAllDBSessions()
	allMutexes := wv.Core.GetAllDBMutexes()

	var mutex *sync.Mutex

	whispersMutex.Lock()
	for i, current := range allSessions {
		if current == nil {
			continue
		}
		mutex = allMutexes[i]
		mutex.Lock()

		var whispers []Whisper
		current.Find(&whispers)

		for _, whisper := range whispers {
			whispersMap[whisper.UniqueId] = &whisper
		}

		mutex.Unlock()
	}
	whispersMutex.Unlock()

	go checkWhispers()
}

func AddWhisper(w *Whisper) {
	s := wv.Core.SessionCollection.GetSession(w.GetDBIndex())
	mutex := wv.Core.SessionCollection.GetMutex(w.GetDBIndex())
	mutex.Lock()
	tx := s.Begin()
	tx.Save(w)
	tx.Commit()
	mutex.Unlock()

	whispersMutex.Lock()
	whispersMap[w.UniqueId] = w
	whispersMutex.Unlock()
}

func GetWhisper(uniqueId string) *Whisper {
	whispersMutex.Lock()
	w := whispersMap[uniqueId]
	whispersMutex.Unlock()
	return w
}

func CreateNewWhisper(result *gotgbot.ChosenInlineResult) *Whisper {
	w := &Whisper{
		Sender: result.From.Id,
		Text:   result.Query,
	}
	w.ParseRecipient()
	w.GenerateUniqueID()
	AddWhisper(w)
	return w
}

func RemoveWhisper(w *Whisper) {
	s := wv.Core.SessionCollection.GetSession(w.GetDBIndex())
	mutex := wv.Core.SessionCollection.GetMutex(w.GetDBIndex())
	mutex.Lock()
	tx := s.Begin()
	tx.Delete(w)
	tx.Commit()
	mutex.Unlock()

	whispersMutex.Lock()
	delete(whispersMap, w.UniqueId)
	whispersMutex.Unlock()
}

// removeWhisperDB will remove the specified whisper ONLY from database.
// this function is an internal function to prevent from deadlock.
func removeWhisperDB(w *Whisper) {
	s := wv.Core.SessionCollection.GetSession(w.GetDBIndex())
	mutex := wv.Core.SessionCollection.GetMutex(w.GetDBIndex())
	mutex.Lock()
	tx := s.Begin()
	tx.Delete(w)
	tx.Commit()
	mutex.Unlock()
}

func checkWhispers() {
	interval := wotoConfig.GetIntervalCheck()
	expiry := wotoConfig.GetExpiry()

	for {
		time.Sleep(interval)
		if whispersMap == nil || whispersMutex == nil {
			return
		}

		whispersMutex.Lock()
		for _, whisper := range whispersMap {
			if whisper.IsExpired(expiry) {
				delete(whispersMap, whisper.UniqueId)

				// prevent from deadlock
				removeWhisperDB(whisper)
			}
		}
		whispersMutex.Unlock()
	}
}
