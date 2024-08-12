/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
