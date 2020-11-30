package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execCmdAndPrint(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println(cmd.String())
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Println(out.String())
	fmt.Println(stderr.String())
	if err != nil {
		panic(err)
	}
}
