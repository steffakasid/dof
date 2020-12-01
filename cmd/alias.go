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
	Short:              "Run any git command on your dot file repository",
	Long: `Basically this command is an alias to run git subcommand.

	In Unix you could think about this command as:
	alias config='/usr/bin/git --git-dir=$HOME/.dof/ --work-tree=$HOME'

  This means this command can be used for any possible git command.
  But note that it's limited to the .dof repository.

	Examples:
	To run 'git --git-dir=$HOME/.dof/ --work-tree=$HOME status'
	you could just run 'dof alias status'

	To view the remote you would just run:
	'dof alias remote -vv'

	To force push to the remote you could just run:
	'dof alias push origin main --force`,
	Run: func(cmd *cobra.Command, args []string) {
		gitAlias.Args = append(gitAlias.Args, args...)
		execCmdAndPrint(gitAlias)
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
