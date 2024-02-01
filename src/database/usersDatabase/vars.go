package usersDatabase

import (
	ws "github.com/ALiwoto/ssg/ssg"
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
