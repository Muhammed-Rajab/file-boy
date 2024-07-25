package codec

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

type EncryptionOp struct {
	FromPath string
	ToPath   string
	Data     []byte
	Salt     []byte
	IV       []byte
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

	combined := []byte{}
	combined = append(combined, []byte("LOVE")...)
	combined = append(combined, eop.Salt...)
	combined = append(combined, eop.IV...)
	combined = append(combined, eop.Data...)

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

func EncryptToZip(zipWriter *zip.Writer, fromPath, relPath string, passphrase []byte) error {
}

func EncryptFromDirToDir(fromPath, toPath string, passphrase []byte) ([]EncryptionOp, error) {

	newZipFile, err := os.Create("./output.zip")
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(newZipFile)

	err = filepath.WalkDir(fromPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(fromPath, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			if relativePath != "." {
				if _, err = zipWriter.Create(relativePath + "/"); err != nil {
					return err
				}
			}
			return nil
		}

		// ! CONTINUE HERE

		return nil
	})

	if err != nil {
		panic(err)
	}

	return nil, nil
}
