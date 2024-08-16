package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Muhammed-Rajab/file-boy/codec"
)

func ExecuteCommandString(execCmd, filePath string, stdin *bytes.Reader, cdc *codec.Codec) error {
	execCmdString := strings.Replace(execCmd, "{1}", filePath, 1)

	cmd := exec.Command("sh", "-c", execCmdString)
	cmd.Stdin = stdin

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if cdc.IsVerbose() {
		if err != nil {
			log.Printf("Error from runinng `%s`: %v\n", execCmdString, err)
			log.Println("continuing")
		}
	}
	fmt.Fprintf(os.Stderr, "[OUT]:\n%s\n", out.String())
	return nil
}
