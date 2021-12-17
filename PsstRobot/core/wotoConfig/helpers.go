package wotoConfig

import (
	"time"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

func ParseConfig(configFile string) (*PsstBotConfig, error) {
	if ConfigSettings != nil {
		return ConfigSettings, nil
	}

	s := &PsstBotConfig{}

	err := strongStringGo.ParseConfig(s, configFile)
	if err != nil {
		return nil, err
	}

	ConfigSettings = s

	return ConfigSettings, nil
}

func LoadConfig() (*PsstBotConfig, error) {
	return ParseConfig("config.ini")
}

func GetCmdPrefixes() []rune {
	return []rune{'/', '!'}
}

func GetBotToken() string {
	if ConfigSettings != nil {
		return ConfigSettings.BotToken
	}
	return ""
}

func IsSingleDb() bool {
	if ConfigSettings != nil {
		return ConfigSettings.SingleDb
	}

	return false
}
func GetIntervalCheck() time.Duration {
	if ConfigSettings != nil {
		return time.Duration(ConfigSettings.IntervalCheck) * time.Minute
	}

	return time.Minute * 5
}

func GetExpiry() time.Duration {
	if ConfigSettings != nil {
		return time.Duration(ConfigSettings.MaxExpiry) * (24 * time.Hour)
	}

	return 7 * (24 * time.Hour)
}
