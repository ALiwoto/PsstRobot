package whisperDatabase

import (
	"sync"
	"time"

	"github.com/ALiwoto/PsstRobot/src/core/wotoConfig"
	wv "github.com/ALiwoto/PsstRobot/src/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"gorm.io/gorm"
)

func LoadAllWhispers() {
	allSessions := wv.Core.GetAllDBSessions()
	allMutexes := wv.Core.GetAllDBMutexes()

	var mutex *sync.Mutex
	expiration := wotoConfig.GetExpiry()

	for i, current := range allSessions {
		if current == nil {
			continue
		}
		mutex = allMutexes[i]
		mutex.Lock()

		var whispers []Whisper
		var expiredWhispers []Whisper
		current.Find(&whispers)

		for currentIndex := 0; currentIndex < len(whispers); currentIndex++ {
			whisper := whispers[currentIndex]
			if whisper.IsExpired(expiration) {
				expiredWhispers = append(expiredWhispers, whisper)
				continue
			}

			whispersMap.Add(whisper.UniqueId, &whisper)
		}

		if len(expiredWhispers) > 0 {
			removeWhispersDBNoLock(expiredWhispers, current)
		}
		mutex.Unlock()
	}

	go checkWhispers()
}

func UpdateWhisper(w *Whisper) {
	index := w.GetDBIndex()
	s := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	mutex.Lock()
	tx := s.Begin()
	tx.Save(w)
	tx.Commit()
	mutex.Unlock()
}

func AddWhisper(w *Whisper) {
	index := w.GetDBIndex()
	s := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	mutex.Lock()
	tx := s.Begin()
	tx.Create(w)
	tx.Commit()
	mutex.Unlock()

	whispersMap.Add(w.UniqueId, w)
}

func GetWhisper(uniqueId string) *Whisper {
	return whispersMap.Get(uniqueId)
}

func CreateNewWhisperFromChosen(result *gotgbot.ChosenInlineResult) *Whisper {
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

// removeWhisperDBNoLock will remove the specified whisper ONLY from database.
// this function is an internal function to prevent from deadlock.
func removeWhispersDBNoLock(w []Whisper, db *gorm.DB) {
	tx := db.Begin()
	for currentIndex := 0; currentIndex < len(w); currentIndex++ {
		tx.Delete(w[currentIndex])
	}
	tx.Commit()
}

func checkWhispers() {
	interval := wotoConfig.GetIntervalCheck()
	expiry := wotoConfig.GetExpiry()
	whispersMap.SetExpiration(expiry)
	whispersMap.SetOnExpiredPtr(func(key string, value *Whisper) {
		removeWhisperDB(value)
	})

	if interval < time.Minute {
		// internal less than 1 minute will fill up the cpu usage
		// it's wisely advised to not use values less than 1 minute.
		interval = time.Minute
	}
	for {
		time.Sleep(interval)
		if whispersMap == nil {
			return
		}

		whispersMap.DoCheck()
	}
}
