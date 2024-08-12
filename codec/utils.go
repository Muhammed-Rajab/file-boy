package codec

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func directoryExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

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
