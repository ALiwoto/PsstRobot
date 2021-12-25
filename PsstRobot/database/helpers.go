package database

import (
	"os"
	"strconv"

	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/logging"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/PsstRobot/PsstRobot/core/wotoValues"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/usersDatabase"
	"github.com/AnimeKaizoku/PsstRobot/PsstRobot/database/whisperDatabase"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartDB() error {
	var err error
	var db *gorm.DB

	single := wotoConfig.IsSingleDb()
	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	wv.Core.GenerateSessionCollection(single)

	if single {
		db, err = gorm.Open(sqlite.Open("psstbot.db"), conf)
		if err != nil {
			return err
		}
		wv.Core.AddDBSession(db)
	} else {
		var current string
		path := "Database" + string(os.PathSeparator)
		for i := 0; i < wv.MultiDbLength; i++ {
			current = "psstbot" + strconv.Itoa(i) + ".db"
			db, err = gorm.Open(sqlite.Open(path+current), conf)
			if err != nil {
				return err
			}

			wv.Core.AddDBSession(db)
		}
	}

	logging.Info("Database connected ")

	//Create tables if they don't exist
	err = wv.Core.AutoMigrateDB(
		whisperDatabase.ModelWhisper,
		usersDatabase.ModelUserHistory,
	)
	if err != nil {
		return err
	}

	logging.Info("Auto-migrated database schema")

	whisperDatabase.LoadAllWhispers()

	return nil
}
