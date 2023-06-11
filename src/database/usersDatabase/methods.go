package usersDatabase

import (
	"time"

	"github.com/AnimeKaizoku/PsstRobot/src/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/src/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

//---------------------------------------------------------

func (m *historyManager) GetUserHistory(ownerId int64) *HistoryCollection {
	h := m.historyMap.Get(ownerId)

	if h != nil && len(h.History) > wotoValues.MaximumHistory {
		h.FixLength()
	}

	return h
}

func (m *historyManager) CreateCollection(ownerId int64) *HistoryCollection {
	h := &HistoryCollection{
		OwnerId: ownerId,
	}

	m.historyMap.Add(ownerId, h)

	return h
}

func (m *historyManager) SetUserHistory(ownerId int64, h *HistoryCollection) {
	m.historyMap.Add(ownerId, h)
}

func (m *historyManager) SetExpiration(expiration time.Duration) {
	m.historyMap.SetExpiration(expiration)
}

func (m *historyManager) cleanUp() {
	m.historyMap.DoCheck()
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
	newHistory := c.History[counter:]
	removed = c.History[:counter]

	c.History = newHistory
	return
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

func (c *HistoryCollection) Clear() {
	c.History = nil
}

//---------------------------------------------------------

func (u *UserData) IsBanned() bool {
	return u.Status == UserStatusBanned
}

func (u *UserData) IsSendingData() bool {
	return u.ChatStatus == UserChatStatusCreating
}

func (u *UserData) IsIdle() bool {
	return u.Status == UserStatusIdle || u.Status == UserStatusBanned
}

func (u *UserData) SetUserStatusToIdle() {
	u.Status = UserStatusIdle
}

func (u *UserData) SetChatStatusToIdle() {
	u.ChatStatus = UserChatStatusIdle
}

//---------------------------------------------------------
