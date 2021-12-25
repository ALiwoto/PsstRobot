package usersDatabase

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GetUserHistory(ownerId int64) *HistoryCollection {
	collection := theManager.GetUserHistory(ownerId)
	if collection != nil {
		if collection.IsEmpty() {
			return nil
		}

		return collection
	}

	collection = getUserHistoryFromDatabase(ownerId)
	if collection == nil {
		theManager.CreateCollection(ownerId)
		// not found in database
		return nil
	} else if collection.IsEmpty() {
		return nil
	}

	theManager.SetUserHistory(ownerId, collection)

	return collection
}

func SaveInHistory(ownerId int64, target *gotgbot.User) {
	if ownerId == target.Id {
		// don't save user itself in the history...
		return
	}

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

	// at this point, we are sure that collection is not nil
	// and that it doesn't contain the target user's history
	history := collection.AddUser(target)

	index := utils.GetDBIndex(ownerId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)

	mutex.Lock()
	tx := session.Begin()
	tx.Create(history)
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

	return &HistoryCollection{
		History: history,
		OwnerId: userId,
	}
}
