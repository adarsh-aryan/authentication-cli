/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"login-sys/auth-client/client"
	"login-sys/auth-client/config"
	"login-sys/shared"

	"github.com/AlecAivazis/survey/v2"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

// enable2faCmd represents the enable2fa command
var enable2faCmd = &cobra.Command{
	Use:   "enable-2fa",
	Args:  cobra.NoArgs,
	Short: "Setup and activate Two-Factor Authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("enable2fa called")

		//load the session config file
		session_config, err := config.Load()

		if err != nil {
			return err
		}

		// request 2FA
		var setup_reply shared.SetUp2FAResponse
		err = client.Client.Call("AuthService.Request2FASetUp", shared.SessionArgs{SessionId: session_config.SessionID}, &setup_reply)

		if err != nil {
			return err
		}
		// render 2fa qrcode and after that ask user to wright the totp code on the cli to verify 2fa and enable it for them.
		renderQrcode(&setup_reply)
		code := promptCodeInput()

		// verify 2fa
		var verify_reply shared.AuthResponse
		err = client.Client.Call("AuthService.Verify2FA", shared.Verify2FArgs{SessionId: session_config.SessionID, Code: code}, &verify_reply)

		if err != nil {
			return err
		}

		fmt.Println(verify_reply.GetMessage())
		return nil
	},
}

func renderQrcode(reply *shared.SetUp2FAResponse) {

	// render an ASCII QR Code directly onto the console window
	fmt.Println("\nScan this QR code with Google Authenticator or Authy:")
	q, err := qrcode.New(reply.URL, qrcode.Medium)
	if err == nil {
		// Print standard terminal-friendly blocks
		fmt.Println(q.ToSmallString(false))
	}

	fmt.Printf("Or enter manual key: %s\n\n", reply.Secret)
}

func promptCodeInput() string {
	// prompt user immediately to verify they linked it correctly
	var code string

	prompt := &survey.Input{Message: "Enter 6 digit code shown in your Autheticator app:"}
	survey.AskOne(prompt, &code, survey.WithValidator(survey.Required))

	return code
}

func init() {
	rootCmd.AddCommand(enable2faCmd)
}
