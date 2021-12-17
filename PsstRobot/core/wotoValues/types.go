package wotoValues

import (
	"sync"

	"gorm.io/gorm"
)

type WotoCore struct {
	SessionCollection *SessionCollection
}

type SessionCollection struct {
	isSingle bool
	// MainSession is active when and only when we are using
	// single mode database.
	MainSession *gorm.DB

	// MainMutex is active when and only when we are using
	// single mode database.
	MainMutex *sync.Mutex

	// SessionMap is active when and only when we are using
	// multi mode database.
	SessionMap map[int]*gorm.DB

	// SessionMutexes is active when and only when we are using
	// multi mode database.
	SessionMutexes map[int]*sync.Mutex
}
