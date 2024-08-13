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
	Short: "encrypt🔒/decrypt🔓 the specified directory to .zip🤐",
	Run: func(cmd *cobra.Command, args []string) {

		// * Get all the flags
		flags := getDirFlags(cmd)
		mode := flags.Mode
		from := flags.From
		to := flags.To
		verbose := flags.Verbose

		validateDirFlags(flags)

		cdc := codec.NewCodec(verbose)

		switch utils.ValidateMode(mode) {
		// ! MAYBE ONE DAY ADD A WAY TO CALL A PROGRAM
		// ! WHICH TAKES IN RELPATH, ENCRYPTED/DECRYPTED DATA
		// ! ETC, FOR EVERY FILE
		// ! OR MAYBE ADD WAY TO OUTPUT THE DATA TO STDOUT
		// ! BUT FOR NOW, THE APP HAS ENOUGH FEATURES FOR ME TO USE IT. Das is genug!
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

type DirFlags struct {
	Verbose bool
	Mode    string
	From    string
	To      string
}

func getDirFlags(cmd *cobra.Command) DirFlags {
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

	return DirFlags{
		Verbose: verbose,
		Mode:    mode,
		From:    from,
		To:      to,
	}
}

func validateDirFlags(flags DirFlags) {
	if exist, err := utils.DirExists(flags.From); !exist {
		log.Fatalf("the directory '%s' does not exists\n", flags.From)
	} else if err != nil {
		log.Fatalln(err)
	}

	if exist, err := utils.DirExists(flags.To); !exist {
		log.Fatalf("the directory '%s' does not exists\n", flags.To)
	} else if err != nil {
		log.Fatalln(err)
	}
}
