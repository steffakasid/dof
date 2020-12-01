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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout <git-repo-url>",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cloneCmd := exec.Command("git", "clone", "--bare", args[0], repoPath)
		execCmdAndPrint(cloneCmd)

		doNotShowUntrackedFiles()

		err := os.Chdir(repoPath)
		doWePanic(err)
		lsCmd := exec.Command("git", "ls-tree", "--name-only", "main")
		filesString := execCmdAndReturn(lsCmd)
		files := strings.Split(filesString, "\n")
		renameOldFiles(files)

		checkoutCmd := *gitAlias
		checkoutCmd.Args = append(checkoutCmd.Args, "checkout", branch)
		execCmdAndPrint(&checkoutCmd)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func renameOldFiles(files []string) {
	for _, file := range files {
		os.Rename(path.Join(workDir, file), path.Join(workDir, file+"_before_dof"))
	}
}
