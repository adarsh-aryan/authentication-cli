/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// exitCmd represents the exit command
var exitCmd = &cobra.Command{
	Use:   "quit",
	Short: "quit program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bye!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(exitCmd)
}
