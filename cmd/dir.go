package cmd

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/Muhammed-Rajab/file-boy/codec"
	"github.com/Muhammed-Rajab/file-boy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "encrypt/decrypt the specified directory to .zip",
	Run: func(cmd *cobra.Command, args []string) {

		// * Get all the flags
		flags := getDirFlags(cmd)
		mode := flags.Mode
		from := flags.From
		to := flags.To
		execCmd := flags.Exec
		verbose := flags.Verbose
		writeToStdout := flags.WriteToStdout

		validateDirFlags(flags)

		cdc := codec.NewCodec(verbose)

		switch ValidateMode(mode) {
		case ENCRYPT:
			passphrase, err := GetPassphraseFromUser(true)
			if err != nil {
				log.Fatalln(err)
			}

			start := time.Now()
			if cdc.IsVerbose() {
				log.Printf("started at %v", start)
			}

			// ! COMMAND EXECUTION
			// command {1}
			// {1}=path in fs
			// Stdin = piped file data
			var zipBuf *bytes.Buffer
			if to != "" {
				zipBuf, err = cdc.EncryptFromDirToZipFile(from, to, passphrase, func(filePath string, eop *codec.EncryptionOp) error {
					if execCmd != "" {
						err := ExecuteCommandString(execCmd, filePath, bytes.NewReader(eop.AsBytes()), &cdc)
						if err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					log.Fatalln(err)
				}
			} else if to == "" {
				zipBuf, err = cdc.EncryptFromDirToZipBuffer(from, passphrase, func(filePath string, eop *codec.EncryptionOp) error {
					if execCmd != "" {
						err := ExecuteCommandString(execCmd, filePath, bytes.NewReader(eop.AsBytes()), &cdc)
						if err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					log.Fatalln(err)
				}
			}

			if writeToStdout {
				_, err = os.Stdout.Write(zipBuf.Bytes())
				if err != nil {
					log.Fatalln(err)
				}
			}

			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully encrypted '%s'. ended at %v, took %d ms.\n", from, end, end.Sub(start).Milliseconds())
			}
		case DECRYPT:
			passphrase, err := GetPassphraseFromUser(false)
			if err != nil {
				log.Fatalln(err)
			}
			start := time.Now()
			if cdc.IsVerbose() {
				log.Printf("started at %v", start)
			}

			// ! COMMAND EXECUTION
			// command {1}
			// {1}=path in fs
			// Stdin=piped file data
			var zipBuf *bytes.Buffer
			if to != "" {
				zipBuf, err = cdc.DecryptFromDirToZipFile(from, to, passphrase, func(filePath string, dop *codec.DecryptionOp) error {
					if execCmd != "" {
						err := ExecuteCommandString(execCmd, filePath, bytes.NewReader(dop.Data), &cdc)
						if err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				zipBuf, err = cdc.DecryptFromDirToZipBuffer(from, passphrase, func(filePath string, dop *codec.DecryptionOp) error {
					if execCmd != "" {
						err := ExecuteCommandString(execCmd, filePath, bytes.NewReader(dop.Data), &cdc)
						if err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					log.Fatalln(err)
				}
			}
			if writeToStdout {
				_, err = os.Stdout.Write(zipBuf.Bytes())
				if err != nil {
					log.Fatalln(err)
				}
			}
			if cdc.IsVerbose() {
				end := time.Now()
				log.Printf("successfully decrypted '%s'. ended at %v, took %d ms.\n", from, end, end.Sub(start).Milliseconds())
			}
		case INVALID:
			log.Fatalf("invalid mode '%s' provided. valid options (are e|E|encrypt|d|D|decrypt)\n", mode)
		}

	},
}

func init() {
	rootCmd.AddCommand(dirCmd)

	dirCmd.PersistentFlags().BoolP("verbose", "v", false, "show detailed ouput")
	viper.BindPFlag("verbose", dirCmd.PersistentFlags().Lookup("verbose"))

	dirCmd.PersistentFlags().BoolP("stdout", "s", false, "writes the encrypted/decrypted data to os.Stdout")
	viper.BindPFlag("stdout", dirCmd.PersistentFlags().Lookup("stdout"))

	dirCmd.PersistentFlags().StringP("from", "f", "", "the path to the directory to encrypt/decrypt from")
	dirCmd.MarkPersistentFlagRequired("from")
	viper.BindPFlag("from", dirCmd.PersistentFlags().Lookup("from"))

	dirCmd.PersistentFlags().StringP("to", "t", "", "the path to the directory to encrypt/decrypt to")
	viper.BindPFlag("to", dirCmd.PersistentFlags().Lookup("to"))

	dirCmd.PersistentFlags().StringP("mode", "m", "e", "the mode(encrypt|eE|decrypt|dD)")
	viper.BindPFlag("mode", dirCmd.PersistentFlags().Lookup("mode"))

	dirCmd.PersistentFlags().StringP("exec", "x", "", "execute the command with path, relative path and encrypted/decrypted data as arguments")
	viper.BindPFlag("exec", dirCmd.PersistentFlags().Lookup("exec"))
}

type DirFlags struct {
	WriteToStdout bool
	Verbose       bool
	Mode          string
	From          string
	To            string
	Exec          string
}

func getDirFlags(cmd *cobra.Command) DirFlags {
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
	exec, err := cmd.PersistentFlags().GetString("exec")
	if err != nil {
		log.Fatalln(err)
	}

	return DirFlags{
		Verbose:       verbose,
		Mode:          mode,
		From:          from,
		To:            to,
		Exec:          exec,
		WriteToStdout: writeToStdOut,
	}
}

func validateDirFlags(flags DirFlags) {
	if exist, err := utils.DirExists(flags.From); !exist {
		log.Fatalf("the directory '%s' does not exists\n", flags.From)
	} else if err != nil {
		log.Fatalln(err)
	}

	// ! if 'to' is provided, then check if path exists. but if 'to', 'stdout' and 'exec' all are not provided, then show error
	if flags.To != "" {
		if exist, err := utils.DirExists(flags.To); !exist {
			log.Fatalf("the directory '%s' does not exists\n", flags.To)
		} else if err != nil {
			log.Fatalln(err)
		}
	} else if flags.To == "" && !flags.WriteToStdout && flags.Exec == "" {
		log.Fatalln("must provide -t <path> -s, -x <command {1}>,otherwise the operation is useless")
	}
}
