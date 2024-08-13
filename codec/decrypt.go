package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"os"
)

type DecryptionOp struct {
	Data []byte
	Salt []byte
	IV   []byte
}

func decrypt(data, passphrase []byte) (*DecryptionOp, error) {
	// prevents from slicing error
	if len(data) < 32 {
		return nil, ErrNotEncryptFile
	}

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
