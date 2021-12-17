package wotoValues

import (
	"sync"

	"gorm.io/gorm"
)

//---------------------------------------------------------

func (w *WotoCore) GenerateSessionCollection(single bool) {
	w.SessionCollection = &SessionCollection{
		isSingle: single,
	}
}

func (w *WotoCore) AddDBSession(session *gorm.DB) {
	w.SessionCollection.AddDBSession(session)
}

func (w *WotoCore) AutoMigrateDB(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}

	return w.SessionCollection.AutoMigrateDB(models...)
}

func (w *WotoCore) GetAllDBSessions() []*gorm.DB {
	return w.SessionCollection.GetAllSessions()
}

func (w *WotoCore) GetAllDBMutexes() []*sync.Mutex {
	return w.SessionCollection.GetAllMutexes()
}

//---------------------------------------------------------

func (c *SessionCollection) AutoMigrateDB(models ...interface{}) error {
	if c.isSingle {
		return c.MainSession.AutoMigrate(models...)
	}

	var err error

	for _, session := range c.SessionMap {
		if session != nil {
			err = session.AutoMigrate(models...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *SessionCollection) AddDBSession(session *gorm.DB) {
	if c.isSingle {
		c.MainSession = session
		return
	}

	l := len(c.SessionMap)
	if l >= MultiDbLength {
		return
	}

	c.SessionMap[l] = session
}

func (c *SessionCollection) GetSession(num int) *gorm.DB {
	if c.isSingle {
		return c.MainSession
	}

	return c.SessionMap[c.CorrectNum(num)]
}

func (c *SessionCollection) GetAllSessions() []*gorm.DB {
	if c.isSingle {
		return []*gorm.DB{c.MainSession}
	}

	var sessions []*gorm.DB
	for _, session := range c.SessionMap {
		sessions = append(sessions, session)
	}

	return sessions
}

func (c *SessionCollection) GetAllMutexes() []*sync.Mutex {
	if c.isSingle {
		return []*sync.Mutex{c.MainMutex}
	}

	var mutexes []*sync.Mutex
	for _, mutex := range c.SessionMutexes {
		mutexes = append(mutexes, mutex)
	}

	return mutexes
}

func (c *SessionCollection) GetMutex(num int) *sync.Mutex {
	if c.isSingle {
		return c.MainMutex
	}

	return c.SessionMutexes[c.CorrectNum(num)]
}

func (c *SessionCollection) GenerateMutexes() {
	if c.isSingle {
		c.MainMutex = &sync.Mutex{}
		return
	}

	c.SessionMutexes = make(map[int]*sync.Mutex, MultiDbLength)
	for i := 0; i < MultiDbLength; i++ {
		c.SessionMutexes[i] = &sync.Mutex{}
	}
}

func (c *SessionCollection) CorrectNum(num int) int {
	if num >= MultiDbFirstIndex && num <= MultiDbLastIndex {
		return num
	}

	return num % MultiDbLength
}

//---------------------------------------------------------
