package codec

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

var ErrNotEncryptFile = errors.New("file not encrypted")
var ErrPathDoesNotExist = errors.New("path doesn't exist")

func generateSalt(length uint) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func deriveKey(passphrase, salt []byte) []byte {
	return pbkdf2.Key(passphrase, salt, 100000, 32, sha256.New)
}
