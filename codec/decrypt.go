package codec

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io/fs"
	"os"
	"path"
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

func DecryptFromDirToDir(fromPath, toPath string, passphrase []byte) ([]DecryptionOp, error) {

	outputZipFile, err := os.Create(path.Join(toPath, "decrypted.zip"))
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(outputZipFile)
	defer zipWriter.Close()

	// Go through all
	err = filepath.WalkDir(fromPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(fromPath, path)
		if err != nil {
			return err
		}

		// if dir, except '.'
		// create a new dir
		if d.IsDir() {
			if relPath != "." {
				if _, err := zipWriter.Create(filepath.Join(relPath + "/")); err != nil {
					return err
				}
			}
			return nil
		}

		// if a file
		if err := addDecryptedFileToZip(zipWriter, path, relPath, passphrase); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func addDecryptedFileToZip(writer *zip.Writer, path, relPath string, passphrase []byte) error {

	// check if the file is of type '.encrypted'
	toDir := filepath.Dir(relPath)
	fileName := filepath.Base(relPath)

	outputFileName := ""
	if strings.HasSuffix(fileName, ".encrypt") {
		outputFileName = fileName[:len(fileName)-8]
	} else {
		outputFileName = fileName
	}

	// ! CHECK IF THE FIRST BYTES ARE 'LOVE'

	dop, err := DecryptFromFile(path, passphrase)
	if err != nil {
		return err
	}

	entry, err := writer.Create(filepath.Join(toDir, outputFileName))
	if err != nil {
		return nil
	}

	_, err = entry.Write(dop.Data)
	if err != nil {
		return nil
	}

	return nil
}
