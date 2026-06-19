package main

import (
	"bufio"
	"fmt"
	"log"
	"login-sys/auth-client/client"
	"login-sys/auth-client/cmd"
	"net/rpc"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	// log.Println("env vars", os.Environ())

	host := os.Getenv("AUTH_SERVER_HOST")
	port := os.Getenv("PORT")

	var err error
	client.Client, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		log.Println("host", host, port)
		log.Fatal("Could not connect to RPC Server. Error:", err)
	}

	defer client.Client.Close()

	// read the input buffer
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("auth-cli > ")

		if err := scanner.Err(); err != nil {
			fmt.Println("scan err", err)
			break
		}

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		input_args := strings.Fields(line)
		cmd.SetInputArgs(input_args) // set root cmd input args

		if err := cmd.Execute(); err != nil {
			fmt.Println("Error:", err)
		}
	}

}
