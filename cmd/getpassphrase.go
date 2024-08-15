package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func GetPassphraseFromUser(confirm bool) ([]byte, error) {
	// ! DEBUG MODE SPECIFIC
	if os.Getenv("UNDER_DEBUG") == "true" {
		log.Println("using passphrase 'pass' under DEBUG mode")
		return []byte("pass"), nil
	}

	fmt.Fprint(os.Stderr, "enter passphrase: ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if confirm {
		fmt.Fprint(os.Stderr, "re-enter passphrase: ")
		reentered, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(passphrase, reentered) {
			return nil, ErrMismatchingPassphrase
		}
	}
	if err != nil {
		return nil, err
	}
	return passphrase, nil
}
