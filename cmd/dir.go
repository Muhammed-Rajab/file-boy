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

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "encrypt or decrypt the specified directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		from, err := cmd.PersistentFlags().GetString("from")
		if err != nil {
			panic(err)
		}
		if exist, err := utils.DirExists(from); !exist {
			panic("directory path does not exist")
		} else if err != nil {
			panic(err)
		}

		to, err := cmd.PersistentFlags().GetString("to")
		if err != nil {
			panic(err)
		}
		if exist, err := utils.DirExists(to); !exist {
			panic("directory path does not exist")
		} else if err != nil {
			panic(err)
		}

		mode, err := cmd.PersistentFlags().GetString("mode")
		if err != nil {
			panic(err)
		}

		switch utils.ValidateMode(mode) {
		case utils.ENCRYPT:
			_, err := codec.EncryptFromDirToZip(from, to, []byte(""))
			if err != nil {
				panic(err)
			}
			fmt.Println("successfully encrypted folder to zip")

		case utils.DECRYPT:
			fmt.Println("decrypt")
		case utils.INVALID:
			panic("invalid mode")
		}

	},
}

func init() {
	rootCmd.AddCommand(dirCmd)

	dirCmd.PersistentFlags().StringP("from", "f", "", "the path to the directory to encrypt/decrypt from")
	dirCmd.MarkPersistentFlagRequired("from")
	viper.BindPFlag("from", dirCmd.PersistentFlags().Lookup("from"))

	dirCmd.PersistentFlags().StringP("to", "t", "", "the path to the directory to encrypt/decrypt to")
	dirCmd.MarkPersistentFlagRequired("to")
	viper.BindPFlag("to", dirCmd.PersistentFlags().Lookup("to"))

	dirCmd.PersistentFlags().StringP("mode", "m", "e", "the mode(encrypt|eE|decrypt|dD)")
	viper.BindPFlag("mode", dirCmd.PersistentFlags().Lookup("mode"))
}
