package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"os"
	"path/filepath"
	"strings"
)

type DecryptionOp struct {
	FromPath string
	ToPath   string
	Data     []byte
	Salt     []byte
	IV       []byte
}

func decrypt(data, passphrase []byte) (*DecryptionOp, error) {
	format := data[:4]
	salt := data[4:20]
	iv := data[20:32]
	encryptedData := data[32:]

	if !bytes.Equal(format, []byte("LOVE")) {
		return nil, ErrNotEncryptFile
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

func DecryptFromFile(filePath string, passphrase []byte) (*DecryptionOp, error) {
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

func DecryptFromToFile(fromPath string, toPath string, passphrase []byte) (*DecryptionOp, error) {
	if !strings.HasSuffix(toPath, "/") {
		toPath += "/"
	}
	toDir := filepath.Dir(toPath)
	fileName := filepath.Base(fromPath)

	dop, err := DecryptFromFile(fromPath, passphrase)
	if err != nil {
		return nil, err
	}

	outputFileName := ""
	if strings.HasSuffix(fileName, ".encrypt") {
		outputFileName = fileName[:len(fileName)-8] + ""
	} else {
		outputFileName = ""
	}

	outputPath := filepath.Join(toDir, outputFileName)
	err = os.WriteFile(outputPath, dop.Data, 0644)
	if err != nil {
		return nil, err
	}

	dop.FromPath = fromPath
	dop.ToPath = outputPath

	return dop, err
}

func DecryptFromZipToDir(fromPath, toPath string, passphrase []byte) ([]DecryptionOp, error) {
	return nil, nil
}
