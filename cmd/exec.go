package cmd

/*
Copyright © 2020 Steffen Rumpf <github@steffen-rumpf.de>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
