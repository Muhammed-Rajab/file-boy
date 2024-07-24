package main

import "github.com/Muhammed-Rajab/file-boy/codec"

func main() {
	passphrase := []byte("password")

	_, err := codec.EncryptFromToFile("./test.jpg", "./encrypted/", passphrase)
	if err != nil {
		panic(err)
	}

	_, err = codec.DecryptFromToFile("./encrypted/test.jpg.encrypt", "./decrypted/", passphrase)
	if err != nil {
		panic(err)
	}
}
