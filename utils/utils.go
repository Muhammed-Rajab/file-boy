package utils

import (
	"bytes"
	"fmt"
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
	fmt.Print("enter passphrase🔒: ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if confirm {
		fmt.Print("re-enter passphrase🔒: ")
		reentered, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
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
