package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

func main() {
	data := []byte("This is my data")
	passphrase := []byte("password")

	encryptedFilePath := "./test-file.encrypt"
	err := encryptToFile(data, passphrase, encryptedFilePath)
	if err != nil {
		panic(err)
	}

	_, err = decryptFromToFile(encryptedFilePath, "./decrypted/", passphrase)
	if err != nil {
		panic(err)
	}
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

type EncryptionOp struct {
	Data []byte
	Salt []byte
	IV   []byte
}

func encrypt(data, passphrase []byte) (*EncryptionOp, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv, err := generateSalt(uint(gcm.NonceSize()))
	if err != nil {
		return nil, err
	}

	encryptedData := gcm.Seal(nil, iv, data, nil)

	// 16 byte salt
	// 32 byte PBKDF2 key
	// AES 256 algo

	return &EncryptionOp{
		Data: encryptedData,
		Salt: salt,
		IV:   iv,
	}, nil
}

func encryptToFile(data, passphrase []byte, filePath string) error {
	op, err := encrypt(data, passphrase)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	_, err = os.Stat(dir)

	if err != nil {
		return err
	}

	combined := []byte{}
	combined = append(combined, op.Salt...)
	combined = append(combined, op.IV...)
	combined = append(combined, op.Data...)

	err = os.WriteFile(filePath, combined, 0644)
	if err != nil {
		return err
	}

	return nil
}

type DecryptionOp struct {
	Data []byte
	Salt []byte
	IV   []byte
}

func decrypt(data, passphrase []byte) (*DecryptionOp, error) {
	// 0:16 -> salt
	// 16:28 -> iv
	// 28: -> encrypted data
	salt := data[:16]
	iv := data[16:28]
	encryptedData := data[28:]

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedData, err := gcm.Open(nil, iv, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return &DecryptionOp{
		Data: decryptedData,
		Salt: salt,
		IV:   iv,
	}, nil
}

func decryptFromFile(filePath string, passphrase []byte) (*DecryptionOp, error) {
	encrypted, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	dop, err := decrypt(encrypted, passphrase)
	if err != nil {
		return nil, err
	}

	return dop, nil
}

func decryptFromToFile(fromPath string, toPath string, passphrase []byte) (*DecryptionOp, error) {
	// take file path
	// extract file name
	// check if to path exists
	// save the encrypted file without the .encrypt to toPath
	toDir := filepath.Dir(toPath)
	fileName := filepath.Base(fromPath)

	dop, err := decryptFromFile(fromPath, passphrase)
	if err != nil {
		return nil, err
	}

	// write the file to path
	outputPath := filepath.Join(toDir, fileName+".decrypted")
	err = os.WriteFile(outputPath, dop.Data, 0644)
	if err != nil {
		return nil, err
	}

	return dop, err
}
