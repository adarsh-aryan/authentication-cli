package main

import (
	"fmt"
	"log"
	"login-sys/auth-server/auth"
	"login-sys/auth-server/db"
	"net"
	"net/rpc"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	//load .env file
	_ = godotenv.Load()

	db := db.AutoMigrateTables()
	auth := &auth.AuthService{
		DB: db,
	}
	err := rpc.Register(auth) // register auth service
	if err != nil {
		log.Fatal("register error:", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))

	if err != nil {
		log.Fatal("Listen error:", err)
	}

	defer listener.Close()
	log.Printf("RPC Server running on PORT %s", os.Getenv("PORT"))
	rpc.Accept(listener)

}
