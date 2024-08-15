package cmd

import "strings"

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
