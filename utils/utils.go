package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func DirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

type OperationMode int

const (
	ENCRYPT OperationMode = iota
	DECRYPT
	INVALID
)

func ValidateMode(mode string) OperationMode {
	if strings.HasPrefix(mode, "encrypt") || mode == "e" || mode == "E" {
		return ENCRYPT
	} else if strings.HasPrefix(mode, "decrypt") || mode == "d" || mode == "D" {
		return DECRYPT
	} else {
		return INVALID
	}
}

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

func StripEncryptFromName(fileName string) string {
	if strings.HasSuffix(fileName, ".encrypt") {
		return fileName[:len(fileName)-8] + ""
	}
	return fileName
}
