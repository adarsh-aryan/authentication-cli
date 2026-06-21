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

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [username] [password]",
	Short: "Log into your account",
	RunE: func(cmd *cobra.Command, args []string) error {

		// get username and password from cmd flags
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}

		otp, err := cmd.Flags().GetString("otp")
		if err != nil {
			return err
		}

		var reply shared.LoginResponse
		err = client.Client.Call("AuthService.Login", shared.LoginArgs{Username: username, Password: password, OTP: otp}, &reply)

		if err != nil {
			return err
		}

		// save the current session in config file
		err = config.Save(reply.SessionId, reply.UserDetails.SessionExpirationTime)
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

		// send message to cli
		fmt.Println(reply.GetMessage())

		return nil
	},
}

func init() {

	// add login command in rootcmd
	rootCmd.AddCommand(loginCmd)

	// add username and password flags
	loginCmd.Flags().StringP("username", "u", "", "username")
	loginCmd.Flags().StringP("password", "p", "", "password")
	loginCmd.Flags().StringP("otp", "o", "", "otp")

	// mark username and password field required
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
