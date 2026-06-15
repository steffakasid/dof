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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout <git-repo-url>",
	Args:  cobra.ExactArgs(1),
	Short: "Clone a bare dotfile repository and checkout tracked files",
	Long: `Clone a bare git repository containing dotfiles and checkout the tracked files to the work tree.

If existing files are tracked by the repository, they are renamed as backups before checkout.

This command:
  1. clones the repository as a bare repo
  2. disables untracked files in the work tree
  3. renames conflicting files to <file>_before_dof
  4. checks out the configured branch

Example:
  dof checkout git@github.com:steffakasid/my-dot-files.git`,
	RunE: func(_ *cobra.Command, args []string) error {
		logger.Info("Cloning bare repo...")
		gitClone := exec.Command("git", "clone", "--bare", args[0], viper.GetString("repository"))
		if err := execCmdAndPrint(gitClone); err != nil {
			return err
		}

		logger.Info("Configure to not show untracked files...")
		if err := doNotShowUntrackedFiles(); err != nil {
			return err
		}

		logger.Info("Rename old files as backup...")
		if err := renameOldFiles(); err != nil {
			return err
		}

		logger.Info("Checkout branch...")
		gitCheckout := *gitAlias
		gitCheckout.Args = append(gitCheckout.Args, "checkout", viper.GetString("branch"))
		if err := execCmdAndPrint(&gitCheckout); err != nil {
			return err
		}

		return applySkipFiles()
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func renameOldFiles() error {
	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", workDir, err)
	}

	gitTree := *gitAlias

	gitTreeArgs := []string{"ls-tree", "--name-only", viper.GetString("branch")}
	gitTree.Args = append(gitTree.Args, gitTreeArgs...)

	filesString, err := execCmdAndReturn(&gitTree)
	if err != nil {
		return err
	}
	files := strings.Split(filesString, "\n")
	skipFiles := viper.GetStringSlice("skip_files")
	for _, file := range files {
		if file == "" {
			continue
		}
		if isSkipped(file, skipFiles) {
			logger.Infof("Skipping %s (in skip_files)", file)
			continue
		}
		logger.Infof("Rename %s to %s", path.Join(workDir, file), path.Join(workDir, file+"_before_dof"))
		if err := os.Rename(path.Join(workDir, file), path.Join(workDir, file+"_before_dof")); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to rename %s: %w", file, err)
		}
	}
	return nil
}
