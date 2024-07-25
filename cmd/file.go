/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Muhammed-Rajab/file-boy/codec"
	"github.com/Muhammed-Rajab/file-boy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "encrypt or decrypt the specified file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// from, existence, permission
		from, err := cmd.PersistentFlags().GetString("from")
		if err != nil {
			panic(err)
		}
		if exist, err := utils.FileExists(from); !exist {
			panic("file path does not exist")
		} else if err != nil {
			panic(err)
		}

		// to, existence, permission
		to, err := cmd.PersistentFlags().GetString("to")
		if err != nil {
			panic(err)
		}
		if exist, err := utils.DirExists(to); !exist {
			panic("dir path does not exist")
		} else if err != nil {
			panic(err)
		}

		// mode, validation
		mode, err := cmd.PersistentFlags().GetString("mode")
		if err != nil {
			panic(err)
		}

		switch utils.ValidateMode(mode) {
		case utils.ENCRYPT:
			// ask for passphrase
			fmt.Print("enter passphrase🔒: ")
			passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				panic(err)
			}
			_, err = codec.EncryptFromToFile(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			fmt.Printf("successfully encrypted\n")
		case utils.DECRYPT:
			// do decryption
			// ask for passphrase
			fmt.Print("enter passphrase🔒: ")
			passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				panic(err)
			}
			_, err = codec.DecryptFromToFile(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			fmt.Printf("successfully decrypted\n")
		case utils.INVALID:
			// throw error as the mode is invalid
			panic("invalid mode")
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.PersistentFlags().StringP("from", "f", "", "the path to the file to encrypt/decrypt from")
	fileCmd.MarkPersistentFlagRequired("from")
	viper.BindPFlag("from", fileCmd.PersistentFlags().Lookup("from"))

	fileCmd.PersistentFlags().StringP("to", "t", "", "the path to the directory to encrypt/decrypt to")
	fileCmd.MarkPersistentFlagRequired("to")
	viper.BindPFlag("to", fileCmd.PersistentFlags().Lookup("to"))

	fileCmd.PersistentFlags().StringP("mode", "m", "e", "the mode(encrypt|eE|decrypt|dD)")
	// fileCmd.MarkPersistentFlagRequired("mode")
	viper.BindPFlag("mode", fileCmd.PersistentFlags().Lookup("mode"))
}
