package usersDatabase

import (
	"sync"
	"time"
)

type UserStatus int

type historyManager struct {
	historyMap   map[int64]*HistoryCollection
	historyMutex *sync.Mutex
}

type HistoryCollection struct {
	History    []UserHistory
	OwnerId    int64
	cachedTime time.Time
}

type UserHistory struct {
	OwnerId    int64  `json:"owner_id"`
	TargetName string `json:"target_name"`
	TargetId   int64  `json:"target_id"`
}

type UserData struct {
	UserId     int64      `json:"user_id" gorm:"primaryKey"`
	Status     UserStatus `json:"status"`
	cachedTime time.Time  `json:"-" gorm:"-" sql:"-"`
}
