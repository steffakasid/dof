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
	"github.com/steffakasid/dof/internal"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout <git-repo-url>",
	Args:  cobra.ExactArgs(1),
	Short: "Checkout a dot file repository from Git and setup dot files",
	Long: `Checkout a dot file repository from Git and setup dot files.

  Note: If you already have dot files which are in the repository they will be renamed as backup!

  This command does the following:
  1. git clone --bare <git-repo-url> <repo-path>
  2. Disable to show untracked files in <work-dir>
  3. rename all files in <work-dir> which are in the dot file repository to e.g. .zshrc_before_dof
  4. git checkout <branch-name>

  Examples:
  dof checkout git@github.com:steffakasid/my-dot-files.git`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := internal.CheckoutDofRepo(workDir, repoFolderName, args[0], viper.GetString("branch"))
		eh.IsFatalError(err)
<<<<<<< HEAD

		doNotShowUntrackedFiles(repo)

		logger.Info("Rename old files as backup...")
		renameOldFiles()

		wt, err := repo.Worktree()
		eh.IsFatalError(err)

		branch, err := repo.Branch(viper.GetString("branch"))

		coOpts := &git.CheckoutOptions{
			Branch: branch.Merge,
			Keep:   true,
		}
		wt.Checkout(coOpts)
=======
>>>>>>> 1e3961e (Refactoring)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
