package usersDatabase

import (
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/utils"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

//---------------------------------------------------------

func (m *historyManager) GetUserHistory(ownerId int64) *HistoryCollection {
	m.historyMutex.Lock()
	h := m.historyMap[ownerId]
	m.historyMutex.Unlock()
	return h
}

func (m *historyManager) CreateCollection(ownerId int64) *HistoryCollection {
	h := &HistoryCollection{
		OwnerId: ownerId,
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

func (c *HistoryCollection) AddUser(user *gotgbot.User) *UserHistory {
	if len(c.History) > wotoValues.MaximumHistory {
		c.History = c.History[1:]
	}

	h := &UserHistory{
		TargetId:   user.Id,
		OwnerId:    c.OwnerId,
		TargetName: utils.GetName(user),
	}

	c.History = append(c.History, *h)

	return h
}
