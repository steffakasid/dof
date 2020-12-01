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

var (
	dontPush bool
	dontPull bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Args:  cobra.NoArgs,
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !dontPush {
			pushCmd := *gitAlias
			pushCmd.Args = append(pushCmd.Args, "push", "origin", branch, "-u")
			execCmdAndPrint(&pushCmd)
		}
		if !dontPull {
			pullCmd := *gitAlias
			pullCmd.Args = append(pullCmd.Args, "pull", "--rebase")
			execCmdAndPrint(&pullCmd)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&dontPull, "push-only", "P", false, "Only push changes to remote repository.")
	syncCmd.Flags().BoolVarP(&dontPush, "pull-only", "p", false, "Only pull changes from remote repository.")
}
