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
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:                "alias",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Short:              "Run a git command against the configured dotfile repository",
	Long: `Alias runs an arbitrary git subcommand against the configured bare repository.

This is equivalent to running git with --git-dir set to the configured repository and --work-tree set to the configured work tree.

Examples:
  dof alias status
  dof alias remote -vv
  dof alias push origin main --force`,
	Run: func(_ *cobra.Command, args []string) {
		gitAlias.Args = append(gitAlias.Args, args...)
		if err := execCmdAndPrint(gitAlias); err != nil {
			logger.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
