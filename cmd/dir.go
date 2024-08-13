package cmd

import (
	"log"
	"time"

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

		if exist, err := utils.DirExists(from); !exist {
			log.Fatalf("the directory '%s' does not exists\n", from)
		} else if err != nil {
			log.Fatalln(err)
		}

		if exist, err := utils.DirExists(to); !exist {
			log.Fatalf("the directory '%s' does not exists\n", to)
		} else if err != nil {
			log.Fatalln(err)
		}

		cdc := codec.NewCodec(verbose)

		switch utils.ValidateMode(mode) {
		// ! MAYBE ONE DAY ADD A WAY TO CALL A PROGRAM
		// ! WHICH TAKES IN RELPATH, ENCRYPTED/DECRYPTED DATA
		// ! ETC, FOR EVERY FILE
		case utils.ENCRYPT:
			passphrase, err := utils.GetPassphraseFromUser(true)
			if err != nil {
				log.Fatalln(err)
			}
			start := time.Now()
			if cdc.IsVerbose() {
				log.Printf("started at %v", start)
			}
			_, err = cdc.EncryptFromDirToZip(from, to, passphrase)
			if err != nil {
				log.Fatalln(err)
			}
			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully encrypted '%s'. ended at %v, took %d ms.\n", from, end, end.Sub(start).Milliseconds())
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
			_, err = cdc.DecryptFromDirToZip(from, to, passphrase)
			if err != nil {
				log.Fatalln(err)
			}
			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully decrypted '%s'. ended at %v, took %d ms.\n", from, end, end.Sub(start).Milliseconds())
			}
		case utils.INVALID:
			log.Fatalf("invalid mode '%s' provided. valid options (are e|E|encrypt|d|D|decrypt)\n", mode)
		}

	},
}

func init() {
	rootCmd.AddCommand(dirCmd)

	dirCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
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
