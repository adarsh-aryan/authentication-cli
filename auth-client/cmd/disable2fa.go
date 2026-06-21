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

// disable2faCmd represents the disable2fa command
var disable2faCmd = &cobra.Command{
	Use:   "disable-2fa",
	Args:  cobra.NoArgs,
	Short: "Disable two-factor authetication",
	RunE: func(cmd *cobra.Command, args []string) error {

		// load the session config file
		session_config, err := config.Load()

		if err != nil {
			return err
		}

		// disable 2fa
		var reply shared.AuthResponse
		err = client.Client.Call("AuthService.Disable2FA", shared.SessionArgs{SessionId: session_config.SessionID}, &reply)

		if err != nil {
			return err
		}

		fmt.Println(reply.GetMessage())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(disable2faCmd)
}
