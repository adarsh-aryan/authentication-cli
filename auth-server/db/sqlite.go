package db

import (
	"login-sys/auth-server/models"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {

	// before connecting to the database , we have to create a data directory where we store auth.db file (for SQL lite DB)
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// create a data directory in current directory
	data_dir := filepath.Join(cwd, "data")
	err = os.MkdirAll(data_dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	db_path := filepath.Join(data_dir, "auth.db")

	return gorm.Open(
		sqlite.Open(db_path),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
}

func AutoMigrateTables() *gorm.DB {

	var err error

	db, err := Connect()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Account{},
	)

	return db
}
