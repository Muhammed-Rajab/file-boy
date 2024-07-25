package utils

import (
	"os"
	"strings"
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
