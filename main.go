/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"login-sys/cmd"
	"login-sys/db"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// migrate the tables in sqlite db (contains users, sessions)
	db.AutoMigrateTables()
	// execute the cli cmd
	cmd.Execute()
}
