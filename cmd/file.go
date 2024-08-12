/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Muhammed-Rajab/file-boy/codec"
	"github.com/Muhammed-Rajab/file-boy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "encrypt or decrypt the specified file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// * Get necessary flags
		verbose, err := cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			panic(err)
		}
		mode, err := cmd.PersistentFlags().GetString("mode")
		if err != nil {
			panic(err)
		}
		from, err := cmd.PersistentFlags().GetString("from")
		if err != nil {
			panic(err)
		}
		to, err := cmd.PersistentFlags().GetString("to")
		if err != nil {
			panic(err)
		}

		if exist, err := utils.FileExists(from); !exist {
			panic("file path does not exist")
		} else if err != nil {
			panic(err)
		}

		if exist, err := utils.DirExists(to); !exist {
			panic("dir path does not exist")
		} else if err != nil {
			panic(err)
		}

		cdc := codec.NewCodec(verbose)

		switch utils.ValidateMode(mode) {
		case utils.ENCRYPT:
			passphrase, err := utils.GetPassphraseFromUser(true)
			if err != nil {
				panic(err)
			}
			_, err = cdc.EncryptFromToFile(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			fmt.Printf("successfully encrypted '%s'\n", from)
		case utils.DECRYPT:
			passphrase, err := utils.GetPassphraseFromUser(false)
			if err != nil {
				panic(err)
			}
			_, err = codec.DecryptFromToFile(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			fmt.Printf("successfully decrypted '%s'\n", from)
		case utils.INVALID:
			panic("invalid mode")
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
	fileCmd.MarkPersistentFlagRequired("verbose")
	viper.BindPFlag("verbose", fileCmd.PersistentFlags().Lookup("verbose"))

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
