package codec

import (
	"archive/zip"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Codec struct {
	verbose bool
}

func NewCodec(verbose bool) Codec {
	return Codec{
		verbose,
	}
}

func (c *Codec) IsVerbose() bool {
	return c.verbose
}

// ENCRYPTION
func (c *Codec) EncryptFromToFile(fromPath, toPath string, passphrase []byte) (*EncryptionOp, error) {
	if !strings.HasSuffix(toPath, "/") {
		toPath += "/"
	}
	toDir := filepath.Dir(toPath)
	fileName := filepath.Base(fromPath)

	eop, err := EncryptFromFile(fromPath, passphrase)
	if err != nil {
		return nil, err
	}

	outputFilePath := filepath.Join(toDir, fileName+".encrypt")

	if exist, err := directoryExists(toDir); !exist {
		return nil, ErrPathDoesNotExist
	} else if err != nil {
		return nil, err
	}

	combined := eop.AsBytes()

	err = os.WriteFile(outputFilePath, combined, 0644)
	if err != nil {
		return nil, err
	}

	return eop, nil
}

func (c *Codec) EncryptFromDirToZip(fromPath, toPath string, passphrase []byte) ([]EncryptionOp, error) {

	// ! MAYBE ADD A WAY TO CHANGE THE NAME OF THE FILE TO
	// ! SOMETHING MORE MEANINGFUL
	newZipFile, err := os.Create(path.Join(toPath, "output.zip"))
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// recursive file tree traversal
	err = filepath.WalkDir(fromPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// relative path fromPath -> path
		relPath, err := filepath.Rel(fromPath, path)
		if err != nil {
			return err
		}

		// directory
		if d.IsDir() {
			if relPath != "." {
				if _, err := zipWriter.Create(relPath + "/"); err != nil {
					return err
				}
			}
			return nil
		}

		// file
		if err := c.writeEncryptedFileToZip(zipWriter, path, relPath, passphrase); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Codec) writeEncryptedFileToZip(writer *zip.Writer, filePath, relPath string, passphrase []byte) error {

	eop, err := EncryptFromFile(filePath, passphrase)
	if err != nil {
		return err
	}

	combined := eop.AsBytes()

	entry, err := writer.Create(relPath + ".encrypt")
	if err != nil {
		return err
	}

	_, err = entry.Write(combined)
	if err != nil {
		return err
	}

	return nil
}

// DECRYPTION
func (c *Codec) DecryptFromToFile(fromPath string, toPath string, passphrase []byte) (*DecryptionOp, error) {
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

	return dop, err
}

func (c *Codec) DecryptFromDirToZip(fromPath, toPath string, passphrase []byte) ([]DecryptionOp, error) {

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

		if d.IsDir() {
			if relPath != "." {
				if _, err := zipWriter.Create(filepath.Join(relPath + "/")); err != nil {
					return err
				}
			}
			return nil
		}

		// if a file
		if err := c.addDecryptedFileToZip(zipWriter, path, relPath, passphrase); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Codec) addDecryptedFileToZip(writer *zip.Writer, path, relPath string, passphrase []byte) error {

	toDir := filepath.Dir(relPath)
	fileName := filepath.Base(relPath)

	outputFileName := ""
	if strings.HasSuffix(fileName, ".encrypt") {
		outputFileName = fileName[:len(fileName)-8]
	} else {
		outputFileName = fileName
	}

	dop, err := DecryptFromFile(path, passphrase)
	// ! DON'T PANIC IF THERE'S A NON ENCRYPTED FILE
	if err == ErrNotEncryptFile && c.verbose {
		log.Println("not encrypted file found")
		return nil
	} else if err != nil {
		return err
	}

	entry, err := writer.Create(filepath.Join(toDir, outputFileName))
	if err != nil {
		return err
	}

	_, err = entry.Write(dop.Data)
	if err != nil {
		return err
	}

	return nil
}
