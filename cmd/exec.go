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
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

// buildGitEnv returns the environment used for every git invocation:
// the inherited OS environment (os.Environ()) with the configured
// git_env entries appended as KEY=value. Values are used verbatim (no
// expansion). When git_env is empty the result equals os.Environ().
func buildGitEnv() []string {
	gitEnv := viper.GetStringMapString("git_env")
	env := os.Environ()
	for k, v := range gitEnv {
		env = append(env, k+"="+v)
	}
	return env
}

// newGitCmd builds an *exec.Cmd for the git binary with the composed
// git_env environment already applied. Used for standalone commands
// (git init --bare, git clone --bare) that cannot copy the gitAlias
// template because the repository does not yet exist.
func newGitCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Env = buildGitEnv()
	return cmd
}

func execCmdAndPrint(cmd *exec.Cmd) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if out.Len() != 0 {
		logger.Info(out.String())
	}
	if stderr.Len() != 0 {
		logger.Error(stderr.String())
	}

	if err != nil {
		return fmt.Errorf("command %s failed: %w", cmd.Path, err)
	}
	return nil
}

func execCmdAndReturn(cmd *exec.Cmd) (string, error) {
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w", cmd.Path, err)
	}
	logger.Info("Output:", string(output))
	return string(output), nil
}

func init() {
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}
