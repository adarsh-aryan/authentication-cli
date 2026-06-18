package db

import (
	"login-sys/auth/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {

	return gorm.Open(
		sqlite.Open("auth.db"),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
}

func AutoMigrateTables() *gorm.DB {

	var err error

	DB, err = Connect()
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Session{},
	)

	return DB
}
