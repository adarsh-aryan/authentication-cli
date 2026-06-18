/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use: "login",
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

		//login user
		err = auth.Login(username, password)
		if err != nil {
			return err
		}
		fmt.Println("login called")
		return nil
	},
}

func init() {

	// add login command in rootcmd
	rootCmd.AddCommand(loginCmd)

	// add username and password flags
	loginCmd.Flags().StringP("username", "u", "", "username")
	loginCmd.Flags().StringP("password", "p", "", "password")

	// mark username and password field required
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}
