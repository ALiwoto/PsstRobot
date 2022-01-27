package usersDatabase

import (
	"time"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func LoadUsersDatabase() {
	go checkUsersData()
}

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

func PrivacyEnabled(userId int64) bool {
	d := GetUserData(userId)
	return d != nil && d.PrivacyMode
}

func HasPrivacy(user *gotgbot.User) bool {
	return PrivacyEnabled(user.Id)
}

func ChangePrivacy(user *gotgbot.User, privacy bool) {
	data := GetUserData(user.Id)
	if data == nil {
		data = &UserData{
			UserId:     user.Id,
			cachedTime: time.Now(),
		}

		userDataMutex.Lock()
		userDataMap[user.Id] = data
		userDataMutex.Unlock()
	}

	if privacy == data.PrivacyMode {
		// prevent from sending unnecessary database queries
		return
	}

	data.PrivacyMode = privacy
	UpdateUserData(data)
}

func IsUserBanned(user *gotgbot.User) bool {
	data := GetUserData(user.Id)
	return data != nil && data.IsBanned()
}

func ChangeUserStatus(user *gotgbot.User, status UserStatus) {
	ChangeUserStatusById(user.Id, status)
}

func ChangeUserStatusById(userId int64, status UserStatus) {
	data := GetUserData(userId)
	if data == nil {
		data = &UserData{
			UserId:     userId,
			cachedTime: time.Now(),
		}

		userDataMutex.Lock()
		userDataMap[userId] = data
		userDataMutex.Unlock()
	}

	if data.Status == status {
		// prevent from sending unnecessary database queries
		return
	}

	data.Status = status
	UpdateUserData(data)
}

func EnablePrivacy(user *gotgbot.User) {
	ChangePrivacy(user, true)
}

func DisablePrivacy(user *gotgbot.User) {
	ChangePrivacy(user, false)
}

func GetUserData(userId int64) *UserData {
	userDataMutex.Lock()
	data := userDataMap[userId]
	userDataMutex.Unlock()
	if data != nil {
		return data
	}

	index := utils.GetDBIndex(userId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	data = &UserData{}

	mutex.Lock()
	session.Model(ModelUserData).Where(
		"user_id = ?", userId,
	).Take(data)
	mutex.Unlock()

	if data.UserId != userId {
		return nil
	}

	userDataMutex.Lock()
	userDataMap[userId] = data
	userDataMutex.Unlock()

	return data
}

func GetUserStatus(user *gotgbot.User) UserStatus {
	data := GetUserData(user.Id)
	if data == nil {
		return UserStatusIdle
	}

	return data.Status
}

func IsUserCreating(user *gotgbot.User) bool {
	return GetUserStatus(user) == UserStatusCreating
}

func UpdateUserData(data *UserData) {
	index := utils.GetDBIndex(data.UserId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)

	mutex.Lock()
	tx := session.Begin()
	tx.Save(data)
	tx.Commit()
	mutex.Unlock()
}

func SaveInHistory(ownerId int64, target *gotgbot.User) {
	if ownerId == target.Id {
		// don't save user itself in the history...
		return
	}

	collection := theManager.GetUserHistory(ownerId)
	if collection != nil {
		if collection.Exists(target.Id) {
			if collection.HasTooMuch() {
				removed := collection.FixLength()
				if len(removed) < 1 {
					return
				}

				index := utils.GetDBIndex(ownerId)
				session := wv.Core.SessionCollection.GetSession(index)
				mutex := wv.Core.SessionCollection.GetMutex(index)

				mutex.Lock()
				session.Delete(removed)
				mutex.Unlock()
			}
			return
		}
	} else {
		// collection doesn't exist?
		// create one
		collection = theManager.CreateCollection(ownerId)
	}

	// at this point, we are sure that collection is not nil
	// and that it doesn't contain the target user's history
	history, removed := collection.AddUser(target)

	index := utils.GetDBIndex(ownerId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)

	mutex.Lock()
	if removed != nil {
		session.Delete(removed)
	}

	tx := session.Begin()
	tx.Create(history)
	tx.Commit()
	mutex.Unlock()
}

func getUserHistoryFromDatabase(userId int64) *HistoryCollection {
	index := utils.GetDBIndex(userId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)
	var history []UserHistoryValue

	mutex.Lock()
	session.Model(ModelUserHistory).Where("owner_id = ?", userId).Find(&history)
	mutex.Unlock()

	return &HistoryCollection{
		History:    history,
		OwnerId:    userId,
		cachedTime: time.Now(),
	}
}

func uIdForUserHistory(ownerId, userId int64) string {
	return utils.ToBase10(ownerId) + "^" + utils.ToBase10(userId)
}

func checkUsersData() {
	interval := wotoConfig.GetIntervalCheck()
	expiry := wotoConfig.GetExpiry()
	for {
		time.Sleep(interval)
		if userDataMap == nil || theManager == nil {
			return
		}

		theManager.cleanUp(expiry)

		userDataMutex.Lock()
		for key, value := range userDataMap {
			if value == nil || value.IsExpired(expiry) {
				delete(userDataMap, key)
			}
		}
		userDataMutex.Unlock()
	}
}
