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
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/steffakasid/dof/internal"
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
		logger.Debugf("Don't push %s", dontPush)
		repo, err := git.PlainOpen(repoPath)
		eh.IsFatalError(err)
		wt, err := repo.Worktree()
		eh.IsFatalError(err)

		if !dontPush {
			status, err := wt.Status()
			eh.IsFatalError(err)
			logger.Info(status.String())
			logger.Debugf("Status of dof %s", status)
			if len(status) > 0 {
				logger.Info("Commiting changed files...")
				opts := &git.CommitOptions{
					All: true,
				}
				_, err = wt.Commit("Synchronized dot files!", opts)
				eh.IsFatalError(err)

				logger.Info("Pushing files")
				pushOpts := &git.PushOptions{RemoteName: "origin", Progress: os.Stdout}
				err = repo.Push(pushOpts)
				eh.IsFatalError(err)
			}
		}
		logger.Debugf("Don't pull %v", dontPull)
		if !dontPull {
			logger.Info("Pulling changes from repo...")
			pullOpts := &git.PullOptions{RemoteName: "origin", Progress: os.Stdout}
			err := wt.Pull(pullOpts)
			eh.IsFatalError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&dontPull, "push-only", "P", false, "Only push changes to remote repository.")
	syncCmd.Flags().BoolVarP(&dontPush, "pull-only", "p", false, "Only pull changes from remote repository.")
	if logger == nil {
		logger = internal.NewOutputLogger(1)
	}
}
