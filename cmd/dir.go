/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/Muhammed-Rajab/file-boy/codec"
	"github.com/Muhammed-Rajab/file-boy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "encrypt or decrypt the specified directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// IF mode is 'e' from is dir
		// IF mode is 'd' from is zip(file)

		// * Get all the flags
		verbose, err := cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			log.Fatalln(err)
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

		if exist, err := utils.DirExists(from); !exist {
			panic("directory path does not exist")
		} else if err != nil {
			panic(err)
		}

		if exist, err := utils.DirExists(to); !exist {
			panic("directory path does not exist")
		} else if err != nil {
			panic(err)
		}

		cdc := codec.NewCodec(verbose)

		switch utils.ValidateMode(mode) {
		case utils.ENCRYPT:
			// If encrypt mode, then check if 'from' dir exists
			// If encrypt mode, then check if 'to' dir exists
			passphrase, err := utils.GetPassphraseFromUser(true)
			if err != nil {
				panic(err)
			}
			_, err = cdc.EncryptFromDirToZip(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			if cdc.IsVerbose() {
				fmt.Println("successfully encrypted folder to zip")
			}

		case utils.DECRYPT:
			passphrase, err := utils.GetPassphraseFromUser(false)
			if err != nil {
				panic(err)
			}
			_, err = cdc.DecryptFromDirToZip(from, to, passphrase)
			if err != nil {
				panic(err)
			}
			if cdc.IsVerbose() {
				fmt.Println("successfully decrypted folder")
			}
		case utils.INVALID:
			panic("invalid mode")
		}

	},
}

func init() {
	rootCmd.AddCommand(dirCmd)

	dirCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
	dirCmd.MarkPersistentFlagRequired("verbose")
	viper.BindPFlag("verbose", dirCmd.PersistentFlags().Lookup("verbose"))

	dirCmd.PersistentFlags().StringP("from", "f", "", "the path to the directory to encrypt/decrypt from")
	dirCmd.MarkPersistentFlagRequired("from")
	viper.BindPFlag("from", dirCmd.PersistentFlags().Lookup("from"))

	dirCmd.PersistentFlags().StringP("to", "t", "", "the path to the directory to encrypt/decrypt to")
	dirCmd.MarkPersistentFlagRequired("to")
	viper.BindPFlag("to", dirCmd.PersistentFlags().Lookup("to"))

	dirCmd.PersistentFlags().StringP("mode", "m", "e", "the mode(encrypt|eE|decrypt|dD)")
	viper.BindPFlag("mode", dirCmd.PersistentFlags().Lookup("mode"))
}
