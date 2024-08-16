package codec

import (
	"archive/zip"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Muhammed-Rajab/file-boy/utils"
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

type EncryptFromDirFn func(filePath string, eop *EncryptionOp) error

func (c *Codec) EncryptFromDirToZip(fromPath, toPath string, passphrase []byte, fn EncryptFromDirFn) ([]EncryptionOp, error) {

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
		eop, err := c.writeEncryptedFileToZip(zipWriter, path, relPath, passphrase)
		if err != nil {
			return err
		}

		err = fn(path, eop)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Codec) writeEncryptedFileToZip(writer *zip.Writer, filePath, relPath string, passphrase []byte) (*EncryptionOp, error) {

	eop, err := EncryptFromFile(filePath, passphrase)
	if err != nil {
		return nil, err
	}

	combined := eop.AsBytes()

	outputFilePath := relPath + ".encrypt"
	entry, err := writer.Create(outputFilePath)
	if err != nil {
		return nil, err
	}

	_, err = entry.Write(combined)
	if err != nil {
		return nil, err
	}

	return eop, nil
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

	outputFileName := utils.StripEncryptFromName(fileName)
	outputPath := filepath.Join(toDir, outputFileName)

	err = os.WriteFile(outputPath, dop.Data, 0644)
	if err != nil {
		return nil, err
	}

	return dop, err
}

type DecryptFromDirFn func(filePath string, dop *DecryptionOp) error

func (c *Codec) DecryptFromDirToZip(fromPath, toPath string, passphrase []byte, fn DecryptFromDirFn) ([]DecryptionOp, error) {

	// ! IMPLEMENT WAY TO DETERMINE A MORE SENSIBLE NAME FOR THE OUTPUT ZIP
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
		dop, err := c.writeDecryptedFileToZip(zipWriter, path, relPath, passphrase)
		if err != nil {
			return err
		}

		err = fn(path, dop)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Codec) writeDecryptedFileToZip(writer *zip.Writer, path, relPath string, passphrase []byte) (*DecryptionOp, error) {

	toDir := filepath.Dir(relPath)
	fileName := filepath.Base(relPath)

	dop, err := DecryptFromFile(path, passphrase)
	if err == ErrNotEncryptFile && c.verbose {
		log.Println("not encrypted file found")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	outputFileName := utils.StripEncryptFromName(fileName)
	entry, err := writer.Create(filepath.Join(toDir, outputFileName))
	if err != nil {
		return nil, err
	}

	_, err = entry.Write(dop.Data)
	if err != nil {
		return nil, err
	}

	return dop, nil
}
