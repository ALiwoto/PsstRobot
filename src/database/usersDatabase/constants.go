package usersDatabase

// users' statuses
const (
	UserStatusIdle   UserStatus = 0
	UserStatusBanned UserStatus = 2
)

const (
	UserChatStatusIdle UserChatStatus = iota
	UserChatStatusCreating
)
