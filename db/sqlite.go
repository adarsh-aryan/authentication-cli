package db

import (
	"login-sys/auth/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {

	dsn := "postgresql://postgres:root@localhost:5432/login_system"

	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
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
