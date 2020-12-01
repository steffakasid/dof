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

func execCmdAndReturn(cmd *exec.Cmd) string {
	output, err := cmd.Output()
	doWePanic(err)
	return string(output)
}

func doWePanic(err error) {
	if err != nil {
		panic(err)
	}
}
