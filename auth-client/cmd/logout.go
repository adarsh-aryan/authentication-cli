/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth-client/client"
	"login-sys/auth-client/config"
	"login-sys/shared"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "end session",
	RunE: func(cmd *cobra.Command, args []string) error {

		// load the session config file
		session_config, err := config.Load()

		if err != nil {
			return err
		}

		var reply shared.AuthResponse
		err = client.Client.Call("AuthService.LogOut", shared.LogoutArgs{SessionId: session_config.SessionID}, &reply)

		if err != nil {
			return err
		}

		// delete the session from the config file
		err = config.Delete()
		if err != nil {
			return err
		}

		// send msg to cli
		fmt.Println(reply.GetMessage())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
