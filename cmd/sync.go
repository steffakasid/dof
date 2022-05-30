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
	"github.com/spf13/viper"
)

var (
	dontPush bool
	dontPull bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Args:  cobra.NoArgs,
	Short: "Synchronize local changes with the remote repository",
	Long: `Synchronize local changes with the remote repository. To do so,
  the command just executes:
  1. git add --all
  2. git commit -a -m "Synchronized dot files"
  3. git push origin <branch-name>
  4. git pull --rebase

  If you pull and there are new files in the remote. You must take care to remove the same files if they're already existing.

  Examples:
  dof sync             - simply push and pull changes to/from the remote repository
  dof sync --push-only - only add, commit and push changes to the remote repository
  dof sync --pull-only - only pull changes from the remote repository`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debugf("Don't push %v", dontPush)
		if !dontPush {
			gitStatus := *gitAlias
			gitStatus.Args = append(gitStatus.Args, "status", "-s")
			status := execCmdAndReturn(&gitStatus)
			logger.Debugf("Status of dof %s", status)
			if len(status) > 0 {
				logger.Info("Commiting changed files...")
				gitCommit := *gitAlias
				gitCommit.Args = append(gitCommit.Args, "commit", "-a", "-m", "Synchronized dot files")
				execCmdAndPrint(&gitCommit)
				logger.Info("Pushing files")
				gitPush := *gitAlias
				gitPush.Args = append(gitPush.Args, "push", "origin", viper.GetString("branch"), "-u")
				execCmdAndPrint(&gitPush)
			}
		}
		logger.Debugf("Don't pull %v", dontPull)
		if !dontPull {
			logger.Info("Pulling changes from repo...")
			gitPull := *gitAlias
			gitPull.Args = append(gitPull.Args, "pull", "--rebase")
			execCmdAndPrint(&gitPull)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&dontPull, "push-only", "P", false, "Only push changes to remote repository.")
	syncCmd.Flags().BoolVarP(&dontPush, "pull-only", "p", false, "Only pull changes from remote repository.")
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}
