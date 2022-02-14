package wotoConfig

type PsstBotConfig struct {
	BotToken      string `section:"main" key:"bot_token"`
	MaxExpiry     int64  `section:"database" key:"max_expiry"`
	SingleDb      bool   `section:"database" key:"single_db"`
	IntervalCheck int    `section:"database" key:"interval_check"`
}
