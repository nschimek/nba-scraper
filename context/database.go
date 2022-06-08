package context

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectToDatabase() *gorm.DB {
	dsn := "connection" // TODO: populate with values from Config
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		Log.Fatal("Could not connect to DB!")
	}

	return db
}
