package usersDatabase

import (
	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
)

// database models
var (
	ModelUserHistory = &UserHistoryValue{}
	ModelUserData    = &UserData{}
)

// caching
var (
	theManager = &historyManager{
		historyMap: ws.NewSafeEMap[int64, HistoryCollection](),
	}

	userDataMap = ws.NewSafeEMap[int64, UserData]()
)
