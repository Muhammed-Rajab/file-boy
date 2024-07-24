package main

import "github.com/Muhammed-Rajab/file-boy/codec"

func main() {
	data := []byte("This is my data")
	passphrase := []byte("password")

	encryptedFilePath := "./test-file.encrypt"
	err := codec.EncryptToFile(data, passphrase, encryptedFilePath)
	if err != nil {
		panic(err)
	}

	_, err = codec.DecryptFromToFile(encryptedFilePath, "./decrypted/", passphrase)
	if err != nil {
		panic(err)
	}
}
