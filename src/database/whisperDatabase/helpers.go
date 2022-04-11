package whisperDatabase

import (
	"sync"
	"time"

	"github.com/AnimeKaizoku/PsstRobot/src/core/wotoConfig"
	wv "github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func LoadAllWhispers() {
	allSessions := wv.Core.GetAllDBSessions()
	allMutexes := wv.Core.GetAllDBMutexes()

	var mutex *sync.Mutex

	for i, current := range allSessions {
		if current == nil {
			continue
		}
		mutex = allMutexes[i]
		mutex.Lock()

		var whispers []Whisper
		current.Find(&whispers)

		for _, whisper := range whispers {
			whispersMap.Add(whisper.UniqueId, &whisper)
		}

		mutex.Unlock()
	}

	go checkWhispers()
}

func AddWhisper(w *Whisper) {
	index := w.GetDBIndex()
	s := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	mutex.Lock()
	tx := s.Begin()
	tx.Save(w)
	tx.Commit()
	mutex.Unlock()

	whispersMap.Add(w.UniqueId, w)
}

func GetWhisper(uniqueId string) *Whisper {
	return whispersMap.Get(uniqueId)
}

func CreateNewWhisper(result *gotgbot.ChosenInlineResult) *Whisper {
	w := &Whisper{
		Sender:          result.From.Id,
		Text:            result.Query,
		InlineMessageId: result.InlineMessageId,
		Type:            WhisperTypePlainText,
	}
	w.ParseRecipient(result)
	w.GenerateUniqueID()
	AddWhisper(w)
	return w
}

func RemoveWhisper(w *Whisper) {
	index := w.GetDBIndex()
	s := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	mutex.Lock()
	tx := s.Begin()
	tx.Delete(w)
	tx.Commit()
	mutex.Unlock()

	whispersMap.Delete(w.UniqueId)
}

// removeWhisperDB will remove the specified whisper ONLY from database.
// this function is an internal function to prevent from deadlock.
func removeWhisperDB(w *Whisper) {
	index := w.GetDBIndex()
	s := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	mutex.Lock()
	tx := s.Begin()
	tx.Delete(w)
	tx.Commit()
	mutex.Unlock()
}

func checkWhispers() {
	interval := wotoConfig.GetIntervalCheck()
	expiry := wotoConfig.GetExpiry()
	whispersMap.SetExpiration(expiry)

	for {
		time.Sleep(interval)
		if whispersMap == nil {
			return
		}

		whispersMap.DoCheck()
	}
}
