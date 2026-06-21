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

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "show current user details",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		//load the session config file
		session_config, err := config.Load()

		if err != nil {
			return err
		}

		var reply shared.AuthResponse
		err = client.Client.Call("AuthService.WhoAmI", shared.SessionArgs{SessionId: session_config.SessionID}, &reply)

		if err != nil {
			return err
		}

		// send msg to cli
		fmt.Println(reply.GetMessage())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
