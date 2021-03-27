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
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/steffakasid/dof/internal"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <file-or-path-to-file>",
	Args:  cobra.ExactArgs(1),
	Short: "Add a new file to the dot file repository",
	Long: `Add a new file to the dot file repository and commit it.

  So basically what it does:
  git add <file>
  git commit -m "Add <file>"

  Examples:
  dof add .zshrc - add .zshrc to dot file repository`,
	Run: func(cmd *cobra.Command, args []string) {
		addAndCommit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	if logger == nil {
		logger = internal.NewOutputLogger(1)
	}
}

func addAndCommit(file string) {
	logger.Infof("Add file %s to dof repository", file)
	// config add .vimrc
	repo, err := git.PlainOpen(repoPath)
	eh.IsFatalError(err)
	wt, err := repo.Worktree()
	eh.IsFatalError(err)
	_, err = wt.Add(file)
	eh.IsFatalError(err)

	opts := &git.CommitOptions{
		All: true,
	}
	wt.Commit(fmt.Sprintf("Add %s", file), opts)
}
