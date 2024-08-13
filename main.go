package main

import (
	"log"
	"os"

	"github.com/Muhammed-Rajab/file-boy/cmd"
)

func main() {
	log.SetOutput(os.Stderr)
	cmd.Execute()
}
