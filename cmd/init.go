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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Args:  cobra.NoArgs,
	Short: "Setup a folder to be used as dot file repository",
	Long: `Setup a folder to be used as dot file repository

  The following steps are executed:
  1. git init --bare <repo-path>
  2. git checkout -B <branch-name>
  3. disable show untracked files in <work-dir>
  4. add .gitignore to <work-dir> to ignore <repo-path>
  5. (optional) add origin remote and set upstream

  Example Usage:
  dof init
  dof init --remote git@github.com:user/dotfiles.git
  dof alias remote add origin <git-repo-url>
  dof add .zshrc
  dof sync --push-only`,
	RunE: func(_ *cobra.Command, _ []string) error {
		logger.Info("Initialize git bare repository...")
		// git init --bare $HOME/.cfg
		gitInit := exec.Command("git", "init", "--bare", viper.GetString("repository"))
		if err := execCmdAndPrint(gitInit); err != nil {
			return err
		}

		logger.Infof("Checkout %s branch\n", viper.GetString("branch"))
		gitCheckout := *gitAlias
		gitCheckout.Args = append(gitCheckout.Args, "checkout", "-B", viper.GetString("branch"))
		if err := execCmdAndPrint(&gitCheckout); err != nil {
			return err
		}

		if err := doNotShowUntrackedFiles(); err != nil {
			return err
		}

		if err := addGitIgnore(); err != nil {
			return err
		}

		if remoteURL != "" {
			if err := addRemote(remoteURL); err != nil {
				return err
			}
		}

		return nil
	},
}

var remoteURL string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&remoteURL, "remote", "", "Remote URL to add as origin after init")
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}

func addRemote(url string) error {
	logger.Infof("Adding remote origin %s...", url)
	gitRemote := *gitAlias
	gitRemote.Args = append(gitRemote.Args, "remote", "add", "origin", url)
	if err := execCmdAndPrint(&gitRemote); err != nil {
		return err
	}

	logger.Infof("Setting upstream to origin/%s...", viper.GetString("branch"))
	gitFetch := *gitAlias
	gitFetch.Args = append(gitFetch.Args, "fetch", "origin")
	// Fetch may fail if remote is empty/unreachable, that's ok
	_ = execCmdAndPrint(&gitFetch)

	gitUpstream := *gitAlias
	gitUpstream.Args = append(gitUpstream.Args, "branch", "--set-upstream-to=origin/"+viper.GetString("branch"))
	// Set upstream may fail if remote branch doesn't exist yet, that's ok
	_ = execCmdAndPrint(&gitUpstream)

	return nil
}

func addGitIgnore() error {
	gitIgnore := path.Join(workDir, ".gitignore")
	file, err := os.Create(gitIgnore)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)

	linesToWrite := []string{repoPathName}
	for _, line := range linesToWrite {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to .gitignore: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush .gitignore: %w", err)
	}

	return addAndCommit(gitIgnore)
}

func doNotShowUntrackedFiles() error {
	// alias config='/usr/bin/git --git-dir=$HOME/.cfg/ --work-tree=$HOME'
	// config config --local status.showUntrackedFiles no
	gitConfigure := *gitAlias
	gitConfigArgs := []string{"config", "--local", "status.showUntrackedFiles", "no"}
	gitConfigure.Args = append(gitConfigure.Args, gitConfigArgs...)
	return execCmdAndPrint(&gitConfigure)
}
