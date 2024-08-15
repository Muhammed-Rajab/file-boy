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
	Short: "encrypt/decrypt the specified file, provided the right passphrase",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// * Get necessary flags
		flags := getFileFlags(cmd)
		from := flags.From
		to := flags.From
		verbose := flags.Verbose
		writeToStdOut := flags.WriteToStdout
		mode := flags.Mode

		// * validate the flags
		validateFileFlags(flags)

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

			// if to path is provided, save to file, else try to stdout it
			var eop *codec.EncryptionOp
			if to != "" {
				eop, err = cdc.EncryptFromToFile(from, to, passphrase)
				if err != nil {
					log.Fatalln(err)
				}
			} else if to == "" {
				eop, err = codec.EncryptFromFile(from, passphrase)
				if err != nil {
					log.Fatalln(err)
				}
			}

			if writeToStdOut {
				_, err = os.Stdout.Write(eop.Data)
				if err != nil {
					log.Fatalln(err)
				}
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

			var dop *codec.DecryptionOp
			if to != "" {
				dop, err = cdc.DecryptFromToFile(from, to, passphrase)
				if err != nil {
					log.Fatalln(err)
				}
			} else if to == "" {
				dop, err = codec.DecryptFromFile(from, passphrase)
				if err != nil {
					log.Fatalln(err)
				}
			}

			if writeToStdOut {
				_, err = os.Stdout.Write(dop.Data)
				if err != nil {
					log.Fatalln(err)
				}
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
	rootCmd.AddCommand(fileCmd)

	fileCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
	viper.BindPFlag("verbose", fileCmd.PersistentFlags().Lookup("verbose"))

	fileCmd.PersistentFlags().BoolP("stdout", "s", false, "writes the encrypted/decrypted data to os.Stdout")
	viper.BindPFlag("stdout", fileCmd.PersistentFlags().Lookup("stdout"))

	fileCmd.PersistentFlags().StringP("from", "f", "", "the path to the file to encrypt/decrypt from")
	fileCmd.MarkPersistentFlagRequired("from")
	viper.BindPFlag("from", fileCmd.PersistentFlags().Lookup("from"))

	fileCmd.PersistentFlags().StringP("to", "t", "", "the path to the directory to encrypt/decrypt to")
	viper.BindPFlag("to", fileCmd.PersistentFlags().Lookup("to"))

	fileCmd.PersistentFlags().StringP("mode", "m", "e", "the mode(encrypt|eE|decrypt|dD)")
	viper.BindPFlag("mode", fileCmd.PersistentFlags().Lookup("mode"))
}

type FileFlags struct {
	WriteToStdout bool
	Verbose       bool
	Mode          string
	From          string
	To            string
}

func getFileFlags(cmd *cobra.Command) FileFlags {
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

	return FileFlags{
		WriteToStdout: writeToStdOut,
		Mode:          mode,
		From:          from,
		To:            to,
		Verbose:       verbose,
	}
}

func validateFileFlags(flags FileFlags) {
	if exist, err := utils.FileExists(flags.From); !exist {
		log.Fatalf("the file path '%s' does not exist\n", flags.From)
	} else if err != nil {
		log.Fatalln(err)
	}

	// ! if 'to' is given, then only validate it, else
	// ! just output the stuff to stdout, if given
	if flags.To != "" {
		if exist, err := utils.DirExists(flags.To); !exist {
			log.Fatalf("the directory path '%s' does not exist\n", flags.To)
		} else if err != nil {
			log.Fatalln(err)
		}
	} else if flags.To == "" && !flags.WriteToStdout {
		log.Fatalln("must provide -t <path> or -s, otherwise the operation is useless")
	}
}
