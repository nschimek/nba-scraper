package context

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Config *Config `Inject:""`
	Gorm   *gorm.DB
}

const (
	dsnFormat = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func createDatabase() *Database {
	return &Database{}
}

func (db *Database) Connect() {
	Log.WithFields(logrus.Fields{"name": db.Config.Database.Name, "location": db.Config.Database.Location}).Info("Connecting to database...")
	dsn := fmt.Sprintf(dsnFormat, db.Config.Database.User, db.Config.Database.Password, db.Config.Database.Location, db.Config.Database.Name)
	gorm, err := gorm.Open(mysql.Open(dsn), db.getGormConfig())

	if err != nil {
		Log.Fatal(err)
	}

	db.Gorm = gorm
}

func (db *Database) getGormConfig() *gorm.Config {
	var logMode logger.LogLevel

	if db.Config.Debug == true {
		logMode = logger.Info
	} else {
		logMode = logger.Error
	}

	return &gorm.Config{Logger: logger.Default.LogMode(logMode)}
}
