package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execCmdAndPrint(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if out.Len() != 0 {
		fmt.Println(out.String())
	}
	if stderr.Len() != 0 {
		fmt.Println(stderr.String())
	}

	doWePanic(err)
}

func doWePanic(err error) {
	if err != nil {
		panic(err)
	}
}
