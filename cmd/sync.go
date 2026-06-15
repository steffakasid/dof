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
	pushOnly bool
	pullOnly bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Args:  cobra.NoArgs,
	Short: "Synchronize local changes with the remote repository",
	Long: `Synchronize local changes with the remote repository.

By default, this command:
  1. checks for local changes
  2. commits changed files
  3. pushes the configured branch to origin
  4. pulls remote changes with rebase

If no local changes exist, the commit and push steps may be skipped.

Examples:
  dof sync             - push and pull changes to/from the remote repository
  dof sync --push-only - only add, commit and push changes to the remote repository
  dof sync --pull-only - only pull changes from the remote repository`,
	RunE: func(_ *cobra.Command, _ []string) error {
		logger.Debugf("Push only: %v", pushOnly)
		if !pullOnly {
			gitStatus := *gitAlias
			gitStatus.Args = append(gitStatus.Args, "status", "-s")
			status, err := execCmdAndReturn(&gitStatus)
			if err != nil {
				return err
			}
			logger.Debugf("Status of dof %s", status)
			needsPush := false
			if len(status) > 0 {
				logger.Info("Committing changed files...")
				gitCommit := *gitAlias
				gitCommit.Args = append(gitCommit.Args, "commit", "-a", "-m", "Synchronized dot files")
				if err := execCmdAndPrint(&gitCommit); err != nil {
					return err
				}
				needsPush = true
			}

			if !needsPush {
				// Check if there are unpushed commits (e.g. from dof add)
				gitLog := *gitAlias
				gitLog.Args = append(gitLog.Args, "log", "@{u}..", "--oneline")
				unpushed, err := execCmdAndReturn(&gitLog)
				if err == nil && len(unpushed) > 0 {
					needsPush = true
				}
			}

			if needsPush {
				logger.Info("Pushing files...")
				gitPush := *gitAlias
				gitPush.Args = append(gitPush.Args, "push", "origin", viper.GetString("branch"), "-u")
				if err := execCmdAndPrint(&gitPush); err != nil {
					return err
				}
			} else {
				logger.Info("Nothing to push.")
			}
		}
		logger.Debugf("Pull only: %v", pullOnly)
		if !pushOnly {
			logger.Info("Pulling changes from repo...")
			gitPull := *gitAlias
			gitPull.Args = append(gitPull.Args, "pull", "--rebase")
			if err := execCmdAndPrint(&gitPull); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&pushOnly, "push-only", "P", false, "Only push changes to remote")
	syncCmd.Flags().BoolVarP(&pullOnly, "pull-only", "p", false, "Only pull changes from remote")
	syncCmd.MarkFlagsMutuallyExclusive("push-only", "pull-only")
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}
