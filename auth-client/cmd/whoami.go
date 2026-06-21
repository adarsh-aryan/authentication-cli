/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
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

		var reply shared.LoginResponse
		err = client.Client.Call("AuthService.WhoAmI", shared.SessionArgs{SessionId: session_config.SessionID}, &reply)

		if err != nil {
			return err
		}

		// show user details to user after login
		data := map[string]any{
			"username":          reply.UserDetails.Username,
			"registration_date": reply.UserDetails.RegistrationDate,
			"mfa_status":        reply.UserDetails.MFAStatus,
			"session_expiry":    reply.UserDetails.SessionExpirationTime,
			"last_login":        reply.UserDetails.LastLoginTime,
		}

		// 2. Marshal to valid JSON bytes
		jsonBytes, _ := json.Marshal(data)

		// 3. Print it perfectly with a newline
		fmt.Println(string(jsonBytes))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
