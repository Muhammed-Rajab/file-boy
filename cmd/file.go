/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"time"

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
		writeToStdOut, err := cmd.PersistentFlags().GetBool("stdout")
		if err != nil {
			log.Fatalln(err)
		}
		verbose, err := cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			log.Fatalln(err)
		}
		mode, err := cmd.PersistentFlags().GetString("mode")
		if err != nil {
			log.Fatalln(err)
		}
		from, err := cmd.PersistentFlags().GetString("from")
		if err != nil {
			log.Fatalln(err)
		}
		to, err := cmd.PersistentFlags().GetString("to")
		if err != nil {
			log.Fatalln(err)
		}

		if exist, err := utils.FileExists(from); !exist {
			log.Fatalf("the file path '%s' does not exist\n", from)
		} else if err != nil {
			log.Fatalln(err)
		}

		if exist, err := utils.DirExists(to); !exist {
			log.Fatalf("the directory path '%s' does not exist\n", to)
		} else if err != nil {
			log.Fatalln(err)
		}

		cdc := codec.NewCodec(verbose)

		switch utils.ValidateMode(mode) {
		case utils.ENCRYPT:
			passphrase, err := utils.GetPassphraseFromUser(true)
			if err != nil {
				log.Fatalln(err)
			}

			start := time.Now()
			if cdc.IsVerbose() {
				log.Printf("started at %v", start)
			}

			eop, err := cdc.EncryptFromToFile(from, to, passphrase)
			if err != nil {
				log.Fatalln(err)
			}

			if writeToStdOut {
				_, err = os.Stdout.Write(eop.Data)
				if err != nil {
					log.Fatalln(err)
				}
			}

			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully encrypted '%s'. ended at %v, took %d seconds.\n", from, end, end.Sub(start).Milliseconds())
			}
		case utils.DECRYPT:
			passphrase, err := utils.GetPassphraseFromUser(false)
			if err != nil {
				log.Fatalln(err)
			}
			start := time.Now()
			if cdc.IsVerbose() {
				log.Printf("started at %v", start)
			}

			dop, err := cdc.DecryptFromToFile(from, to, passphrase)
			if err != nil {
				log.Fatalln(err)
			}

			if writeToStdOut {
				_, err = os.Stdout.Write(dop.Data)
				if err != nil {
					log.Fatalln(err)
				}
			}

			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully decrypted '%s'. ended at %v, took %d seconds.\n", from, end, end.Sub(start).Milliseconds())
			}
		case utils.INVALID:
			log.Fatalf("invalid mode '%s' provided. valid options (are e|E|encrypt|d|D|decrypt)\n", mode)
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
	viper.BindPFlag("verbose", fileCmd.PersistentFlags().Lookup("verbose"))

	fileCmd.PersistentFlags().BoolP("stdout", "s", false, "writes the encrypted/decrypted data to os.Stdout")
	viper.BindPFlag("stdout", fileCmd.PersistentFlags().Lookup("stdout"))

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
