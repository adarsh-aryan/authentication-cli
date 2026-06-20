package main

import (
	"fmt"
	"io"
	"log"
	"login-sys/auth-client/client"
	"login-sys/auth-client/cmd"
	"login-sys/auth-client/utils"
	"net/rpc"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/term"
)

// rwWrapper bundles Stdin and Stdout into a single io.ReadWriter
type rwWrapper struct {
	io.Reader
	io.Writer
}

func main() {

	_ = godotenv.Load()

	host := os.Getenv("AUTH_SERVER_HOST")
	port := os.Getenv("PORT")

	var err error
	client.Client, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		log.Println("host", host, port)
		log.Fatal("Could not connect to RPC Server. Error:", err)
	}

	defer client.Client.Close()

	// 1. Check if we are actually in a terminal (character device)
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		log.Fatal("Current environment is not a terminal")
	}

	// 2. Put terminal into raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal("Failed to enable raw mode:", err)
	}

	// Create a safe restore function
	restoreTerm := func() {
		_ = term.Restore(fd, oldState)
	}
	defer restoreTerm()

	// 3. Pass both Stdin and Stdout to the terminal controller
	t := term.NewTerminal(rwWrapper{os.Stdin, os.Stdout}, "auth-cli > ")

	// register autocomplete callback hook here
	t.AutoCompleteCallback = utils.AutoCompleteHook

	// create the REPL loop

	for {

		// readline natively handles up/down arrow keys for history!
		line, err := t.ReadLine()
		if err != nil {
			// handles Ctrl+C or Ctrl+D graceful exits
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		input_args := strings.Fields(line)

		// temporarily restore terminal store before running the cmd execution
		// this prevents command outputs from formatting strangely in raw mode
		term.Restore(int(os.Stdin.Fd()), oldState)

		cmd.SetInputArgs(input_args)
		if err := cmd.Execute(); err != nil {
			fmt.Println("Error: ", err)
		}

		// reenter raw mode for next prompt iteration
		oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}
	}

}
