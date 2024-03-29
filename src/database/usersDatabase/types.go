package usersDatabase

import (
	ws "github.com/ALiwoto/ssg/ssg"
)

type (
	UserStatus     int
	UserChatStatus int
)

type historyManager struct {
	historyMap *ws.SafeEMap[int64, HistoryCollection]
}

type HistoryCollection struct {
	History []UserHistoryValue
	OwnerId int64
}

type UserHistoryValue struct {
	UniqueId   string `json:"unique_id" gorm:"primaryKey"`
	OwnerId    int64  `json:"owner_id"`
	TargetName string `json:"target_name"`
	TargetId   int64  `json:"target_id"`
}

type UserData struct {
	UserId            int64      `json:"user_id" gorm:"primaryKey"`
	Status            UserStatus `json:"user_status"`
	PrivacyMode       bool       `json:"privacy_mode"`
	IsHistoryDisabled bool       `json:"is_history_disabled"`

	ChatStatus UserChatStatus `json:"-" gorm:"-" sql:"-"`
}
