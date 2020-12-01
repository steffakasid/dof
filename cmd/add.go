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
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		addAndCommit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addAndCommit(file string) {
	// config add .vimrc
	gitAdd := *gitAlias
	gitAddArgs := []string{"add", file}
	gitAdd.Args = append(gitAdd.Args, gitAddArgs...)
	execCmdAndPrint(&gitAdd)

	// config commit -m "Add vimrc"
	gitCommit := *gitAlias
	gitCommitArgs := []string{"commit", "-m", "Add " + file}
	gitCommit.Args = append(gitCommit.Args, gitCommitArgs...)
	execCmdAndPrint(&gitCommit)
}
