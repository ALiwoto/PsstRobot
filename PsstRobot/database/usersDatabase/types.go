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
	History    []UserHistoryValue
	OwnerId    int64
	cachedTime time.Time
}

type UserHistoryValue struct {
	UniqueId   string `json:"unique_id" gorm:"primaryKey"`
	OwnerId    int64  `json:"owner_id"`
	TargetName string `json:"target_name"`
	TargetId   int64  `json:"target_id"`
}

type UserData struct {
	UserId     int64      `json:"user_id" gorm:"primaryKey"`
	Status     UserStatus `json:"status"`
	cachedTime time.Time  `json:"-" gorm:"-" sql:"-"`
}
