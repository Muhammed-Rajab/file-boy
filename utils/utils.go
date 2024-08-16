package utils

import (
	"os"
	"strings"
)

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func DirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func StripEncryptFromName(fileName string) string {
	if strings.HasSuffix(fileName, ".encrypt") {
		return fileName[:len(fileName)-8] + ""
	}
	return fileName
}
