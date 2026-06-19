/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth-client/client"
	"login-sys/shared"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register [username] [password]",
	Short: "Register your account",
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

		// register user
		var reply shared.AuthResponse
		err = client.Client.Call("AuthService.Register", shared.RegisterArgs{Username: username, Password: password}, &reply)

		if err != nil {
			return err
		}

		// send msg to cli
		fmt.Println(reply.GetMessage())

		return nil
	},
}

func init() {

	// add register command in rootcmd
	rootCmd.AddCommand(registerCmd)

	// add username and password flags
	registerCmd.Flags().StringP("username", "u", "", "username")
	registerCmd.Flags().StringP("password", "p", "", "password")

	// mark username and password field required
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
}
