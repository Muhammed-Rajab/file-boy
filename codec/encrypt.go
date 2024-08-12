package codec

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"os"
)

type EncryptionOp struct {
	FromPath string
	ToPath   string
	Data     []byte
	Salt     []byte
	IV       []byte
}

func (eop *EncryptionOp) AsBytes() []byte {
	combined := []byte{}
	combined = append(combined, []byte("LOVE")...)
	combined = append(combined, eop.Salt...)
	combined = append(combined, eop.IV...)
	combined = append(combined, eop.Data...)
	return combined
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

func encryptFromFile(filePath string, passphrase []byte) (*EncryptionOp, error) {
	plain, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	eop, err := encrypt(plain, passphrase)
	if err != nil {
		return nil, err
	}

	return eop, nil
}

func addEncryptedFileToZip(zipWriter *zip.Writer, filePath, relPath string, passphrase []byte) error {

	eop, err := encryptFromFile(filePath, passphrase)
	if err != nil {
		return err
	}

	combined := eop.AsBytes()

	zipFileEntry, err := zipWriter.Create(relPath + ".encrypt")
	if err != nil {
		return err
	}

	_, err = zipFileEntry.Write(combined)
	if err != nil {
		return err
	}

	return nil
}
