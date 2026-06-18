/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:  "whoami",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		username, err := auth.WhoAmI()
		if err != nil {
			return err
		}

		fmt.Println(username)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
