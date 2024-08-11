package codec

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func EncryptFromToFile(fromPath, toPath string, passphrase []byte) (*EncryptionOp, error) {
	if !strings.HasSuffix(toPath, "/") {
		toPath += "/"
	}
	toDir := filepath.Dir(toPath)
	fileName := filepath.Base(fromPath)

	eop, err := EncryptFromFile(fromPath, passphrase)
	if err != nil {
		return nil, err
	}

	outputPath := filepath.Join(toDir, fileName+".encrypt")

	if exist, err := directoryExists(toDir); !exist {
		return nil, ErrPathDoesNotExist
	} else if err != nil {
		return nil, err
	}

	combined := eop.AsBytes()

	err = os.WriteFile(outputPath, combined, 0644)
	if err != nil {
		return nil, err
	}

	eop.FromPath = fromPath
	eop.ToPath = outputPath

	return eop, nil
}

func EncryptFromFile(filePath string, passphrase []byte) (*EncryptionOp, error) {
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

// func EncryptToZip(zipWriter *zip.Writer, fromPath, relPath string, passphrase []byte) error {
// }

func EncryptFromDirToZip(fromPath, toPath string, passphrase []byte) ([]EncryptionOp, error) {

	// Check if the to path exists and all
	newZipFile, err := os.Create(path.Join(toPath, "output.zip"))
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	err = filepath.WalkDir(fromPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// relative to from path
		relPath, err := filepath.Rel(fromPath, path)
		if err != nil {
			return err
		}

		// If it's a directory
		if d.IsDir() {
			if relPath != "." {
				if _, err := zipWriter.Create(relPath + "/"); err != nil {
					return err
				}
			}
			return nil
		}

		// If it's a file
		if err := addEncryptedFileToZip(zipWriter, path, relPath, passphrase); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func addEncryptedFileToZip(zipWriter *zip.Writer, filePath, relPath string, passphrase []byte) error {

	eop, err := EncryptFromFile(filePath, passphrase)
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
