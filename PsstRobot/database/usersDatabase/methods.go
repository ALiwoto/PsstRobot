package usersDatabase

import (
	"time"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

//---------------------------------------------------------

func (m *historyManager) GetUserHistory(ownerId int64) *HistoryCollection {
	m.historyMutex.Lock()
	h := m.historyMap[ownerId]
	m.historyMutex.Unlock()
	if len(h.History) > wotoValues.MaximumHistory {
		h.FixLength()
	}

	return h
}

func (m *historyManager) CreateCollection(ownerId int64) *HistoryCollection {
	h := &HistoryCollection{
		OwnerId:    ownerId,
		cachedTime: time.Now(),
	}

	m.historyMutex.Lock()
	m.historyMap[ownerId] = h
	m.historyMutex.Unlock()

	return h
}

func (m *historyManager) SetUserHistory(ownerId int64, h *HistoryCollection) {
	m.historyMutex.Lock()
	m.historyMap[ownerId] = h
	m.historyMutex.Unlock()
}

func (m *historyManager) cleanUp(expiry time.Duration) {
	m.historyMutex.Lock()

	for k, v := range m.historyMap {
		if v == nil || v.IsExpired(expiry) {
			delete(m.historyMap, k)
		}
	}

	m.historyMutex.Unlock()
}

//---------------------------------------------------------

func (c *HistoryCollection) Exists(targetId int64) bool {
	if c.IsEmpty() {
		return false
	}

	for _, h := range c.History {
		if h.TargetId == targetId {
			return true
		}
	}

	return false
}

func (c *HistoryCollection) IsEmpty() bool {
	return len(c.History) == 0
}

func (c *HistoryCollection) HasTooMuch() bool {
	return len(c.History) > wotoValues.MaximumHistory
}

func (c *HistoryCollection) FixLength() (removed []UserHistoryValue) {
	counter := len(c.History) - wotoValues.MaximumHistory
	c.History = c.History[counter:]
	return c.History[:counter]
}

func (c *HistoryCollection) AddUser(user *gotgbot.User) (new, removed *UserHistoryValue) {
	if len(c.History) > wotoValues.MaximumHistory {
		removed = &c.History[0]
		c.History = c.History[1:]
	}

	h := &UserHistoryValue{
		UniqueId:   uIdForUserHistory(c.OwnerId, user.Id),
		TargetId:   user.Id,
		OwnerId:    c.OwnerId,
		TargetName: utils.GetName(user),
	}

	c.History = append(c.History, *h)

	new = h

	return
}

func (c *HistoryCollection) IsExpired(expiry time.Duration) bool {
	return time.Since(c.cachedTime) > expiry
}

//---------------------------------------------------------

func (u *UserData) IsExpired(expiry time.Duration) bool {
	return time.Since(u.cachedTime) > expiry
}

func (u *UserData) IsBanned() bool {
	return u.Status == UserStatusBanned
}

func (u *UserData) IsSendingData() bool {
	return u.Status == UserStatusCreating
}

func (u *UserData) IsIdle() bool {
	return u.Status == UserStatusIdle || u.Status == UserStatusBanned
}

func (u *UserData) SetToIdle() {
	u.Status = UserStatusIdle
}

//---------------------------------------------------------
