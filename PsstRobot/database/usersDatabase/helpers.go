package usersDatabase

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GetUserHistory(userId int64) *HistoryCollection {
	h := theManager.GetUserHistory(userId)
	if h != nil {
		return h
	}

	h = getUserHistoryFromDatabase(userId)
	if h == nil {
		// not found in database
		return nil
	}

	theManager.SetUserHistory(userId, h)

	return h
}

func SaveInHistory(ownerId int64, target *gotgbot.User) {
	collection := theManager.GetUserHistory(ownerId)
	if collection != nil {
		if collection.Exists(target.Id) {
			return
		}
	} else {
		// collection doesn't exist?
		// create one
		collection = theManager.CreateCollection(ownerId)
	}

	h := collection.AddUser(target)

	index := utils.GetDBIndex(ownerId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)

	mutex.Lock()
	tx := session.Begin()
	tx.Create(h)
	tx.Commit()
	mutex.Unlock()
}

func getUserHistoryFromDatabase(userId int64) *HistoryCollection {
	index := utils.GetDBIndex(userId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	var history []UserHistory

	mutex.Lock()
	session.Model(ModelUserHistory).Where("owner_id = ?", userId).Find(&history)
	mutex.Unlock()

	return nil
}
