package codec

import (
	"archive/zip"
	"io/fs"
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

// ENCRYPTION
func (c *Codec) EncryptFromDirToZip(fromPath, toPath string, passphrase []byte) ([]EncryptionOp, error) {

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

func (c *Codec) EncryptFromToFile(fromPath, toPath string, passphrase []byte) (*EncryptionOp, error) {
	if !strings.HasSuffix(toPath, "/") {
		toPath += "/"
	}
	toDir := filepath.Dir(toPath)
	fileName := filepath.Base(fromPath)

	eop, err := encryptFromFile(fromPath, passphrase)
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

// DECRYPTION
