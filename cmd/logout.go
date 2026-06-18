/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use: "logout",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := auth.LogOut()
		if err != nil {
			return err
		}
		fmt.Println("logout called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
