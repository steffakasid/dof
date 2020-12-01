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
package cmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Args:  cobra.NoArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// git init --bare $HOME/.cfg
		gitInit := exec.Command("git", "init", "--bare", repoPath)

		execCmdAndPrint(gitInit)

		// alias config='/usr/bin/git --git-dir=$HOME/.cfg/ --work-tree=$HOME'
		// config config --local status.showUntrackedFiles no
		gitConfigure := gitAlias
		gitConfigArgs := []string{"config", "--local", "status.showUntrackedFiles", "no"}
		gitConfigure.Args = append(gitConfigure.Args, gitConfigArgs...)
		execCmdAndPrint(gitConfigure)

		addGitIgnore()

		err := viper.SafeWriteConfig()
		if err != nil {
			log.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func addGitIgnore() {

	file, err := os.Create(path.Join(workDir, ".gitignore"))
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	linesToWrite := []string{repoPathName}
	for _, line := range linesToWrite {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	writer.Flush()
}
