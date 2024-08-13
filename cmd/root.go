/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file-boy",
	Short: "a no-bs encryption/decryption cli, made with 💖",
	Long: `a no-bs encryption/decryption cli, made with 💖.

file:
	- encrypt
	- decrypt
	- output encrypted data to stdout
	- output decrypted data to stdout
dir:
	- encrypt directory to zip
	- decrypt directory of encrypted files to zip
`,
	Example: "file-boy file <args>\nfile-boy dir <args>",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
