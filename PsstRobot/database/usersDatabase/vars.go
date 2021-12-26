package usersDatabase

import "sync"

// database models
var (
	ModelUserHistory = &UserHistory{}
	ModelUserData    = &UserData{}
)

// caching
var (
	theManager = &historyManager{
		historyMutex: &sync.Mutex{},
		historyMap:   make(map[int64]*HistoryCollection),
	}

	userDataMap   = make(map[int64]*UserData)
	userDataMutex = &sync.Mutex{}
)
