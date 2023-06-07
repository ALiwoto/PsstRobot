package usersDatabase

import (
	"time"

	"github.com/AnimeKaizoku/PsstRobot/src/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/src/core/wotoConfig"
	wv "github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	ws "github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func LoadUsersDatabase() {
	go checkUsersData()
}

func GetUserHistory(ownerId int64) *HistoryCollection {
	if IsHistoryDisabled(ownerId) {
		return nil
	}

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

func ChangeUserWhisperHistory(userId int64, disabled bool) {
	data := GetUserData(userId)
	if data == nil {
		data = &UserData{
			UserId: userId,
		}

		userDataMap.Add(userId, data)
	}

	if data.IsHistoryDisabled == disabled {
		// prevent from sending unnecessary database queries
		return
	}

	data.IsHistoryDisabled = disabled
	UpdateUserData(data)
}

func ChangePrivacy(user *gotgbot.User, privacy bool) {
	data := GetUserData(user.Id)
	if data == nil {
		data = &UserData{
			UserId: user.Id,
		}

		userDataMap.Add(user.Id, data)
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

func IsHistoryDisabled(userId int64) bool {
	data := GetUserData(userId)
	return data != nil && data.IsHistoryDisabled
}

func ChangeUserStatus(user *gotgbot.User, status UserStatus) {
	ChangeUserStatusById(user.Id, status)
}

func ChangeUserChatStatus(user *gotgbot.User, status UserChatStatus) {
	ChangeUserChatStatusById(user.Id, status)
}

func ChangeUserStatusById(userId int64, status UserStatus) {
	data := GetUserData(userId)
	if data == nil {
		data = &UserData{
			UserId: userId,
		}

		userDataMap.Add(userId, data)
	}

	if data.Status == status {
		// prevent from sending unnecessary database queries
		return
	}

	data.Status = status
	UpdateUserData(data)
}

func ChangeUserChatStatusById(userId int64, status UserChatStatus) {
	data := GetUserData(userId)
	if data == nil {
		data = &UserData{
			UserId: userId,
		}

		userDataMap.Add(userId, data)
	}

	if data.ChatStatus == status {
		// prevent from sending unnecessary database queries
		return
	}

	data.ChatStatus = status
	UpdateUserData(data)
}

func EnablePrivacy(user *gotgbot.User) {
	ChangePrivacy(user, true)
}

func DisablePrivacy(user *gotgbot.User) {
	ChangePrivacy(user, false)
}

func EnableUserWhisperHistory(user *gotgbot.User) {
	ChangeUserWhisperHistory(user.Id, false)
}

func DisableUserWhisperHistory(user *gotgbot.User) {
	ChangeUserWhisperHistory(user.Id, true)
}

func ClearUserWhisperHistory(userId int64) bool {
	collection := GetUserHistory(userId)
	if collection == nil || len(collection.History) == 0 {
		return false
	}

	collection.Clear()

	index := utils.GetDBIndex(userId)
	session := wv.Core.SessionCollection.GetSession(index)
	mutex := wv.Core.SessionCollection.GetMutex(index)

	mutex.Lock()
	session.Model(ModelUserHistory).Delete(
		"owner_id = ?", userId,
	)
	mutex.Unlock()

	return true
}

func GetUserData(userId int64) *UserData {
	data := userDataMap.Get(userId)
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

	userDataMap.Add(userId, data)

	return data
}

func GetUserStatus(user *gotgbot.User) UserChatStatus {
	data := GetUserData(user.Id)
	if data == nil {
		return UserChatStatusIdle
	}

	return data.ChatStatus
}

func IsUserCreating(user *gotgbot.User) bool {
	return GetUserStatus(user) == UserChatStatusCreating
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
	if ownerId == target.Id || IsHistoryDisabled(ownerId) {
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
		History: history,
		OwnerId: userId,
	}
}

func uIdForUserHistory(ownerId, userId int64) string {
	return ws.ToBase10(ownerId) + "^" + ws.ToBase10(userId)
}

func checkUsersData() {
	interval := wotoConfig.GetIntervalCheck()
	expiry := wotoConfig.GetExpiry()
	theManager.SetExpiration(expiry)
	userDataMap.SetExpiration(expiry)
	for {
		time.Sleep(interval)
		if userDataMap == nil || theManager == nil {
			return
		}

		theManager.cleanUp()
		userDataMap.DoCheck()
	}
}
