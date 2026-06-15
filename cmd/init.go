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
	Short: "Initialize a bare dotfile repository",
	Long: `Initialize a bare git repository for dotfiles.

This command:
  1. creates the bare repository at the configured path
  2. checks out the configured branch
  3. disables untracked files in the work tree
  4. adds .gitignore to ignore the repository directory
  5. optionally adds an origin remote and sets upstream

Example usage:
  dof init
  dof init --remote git@github.com:user/dotfiles.git
  dof add .zshrc
  dof sync --push-only`,
	RunE: func(_ *cobra.Command, _ []string) error {
		logger.Info("Initialize git bare repository...")
		// git init --bare $HOME/.cfg
		gitInit := exec.Command("git", "init", "--bare", viper.GetString("repository"))
		if err := execCmdAndPrint(gitInit); err != nil {
			return err
		}

		if err := ensureLocalGitIdentity(); err != nil {
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

		return applySkipFiles()
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

func ensureLocalGitIdentity() error {
	if err := ensureLocalGitConfigValue("user.name", "dof test"); err != nil {
		return err
	}
	return ensureLocalGitConfigValue("user.email", "dof-test@example.com")
}

func ensureLocalGitConfigValue(key, fallback string) error {
	gitGet := *gitAlias
	gitGet.Args = append(gitGet.Args, "config", "--get", key)
	if _, err := gitGet.Output(); err == nil {
		return nil
	}

	gitSet := *gitAlias
	gitSet.Args = append(gitSet.Args, "config", "--local", key, fallback)
	return execCmdAndPrint(&gitSet)
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

func applySkipFiles() error {
	skipFiles := viper.GetStringSlice("skip_files")
	if len(skipFiles) == 0 {
		return nil
	}

	logger.Info("Configuring sparse-checkout to skip files...")

	// Enable sparse-checkout in non-cone mode
	gitConfig := *gitAlias
	gitConfig.Args = append(gitConfig.Args, "config", "--local", "core.sparseCheckout", "true")
	if err := execCmdAndPrint(&gitConfig); err != nil {
		return fmt.Errorf("failed to enable sparse-checkout: %w", err)
	}

	// Write sparse-checkout patterns: include all, then exclude each skip file
	sparseFile := path.Join(viper.GetString("repository"), "info", "sparse-checkout")
	if err := os.MkdirAll(path.Dir(sparseFile), 0700); err != nil {
		return fmt.Errorf("failed to create info directory: %w", err)
	}

	file, err := os.Create(sparseFile)
	if err != nil {
		return fmt.Errorf("failed to create sparse-checkout file: %w", err)
	}
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString("/*\n"); err != nil {
		return fmt.Errorf("failed to write sparse-checkout pattern: %w", err)
	}
	for _, skip := range skipFiles {
		if _, err := writer.WriteString("!" + skip + "\n"); err != nil {
			return fmt.Errorf("failed to write sparse-checkout pattern: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush sparse-checkout file: %w", err)
	}

	// Apply sparse-checkout to working tree
	gitReadTree := *gitAlias
	gitReadTree.Args = append(gitReadTree.Args, "read-tree", "-mu", "HEAD")
	if err := execCmdAndPrint(&gitReadTree); err != nil {
		return fmt.Errorf("failed to apply sparse-checkout: %w", err)
	}

	return nil
}

func isSkipped(file string, skipFiles []string) bool {
	for _, skip := range skipFiles {
		if file == skip {
			return true
		}
	}
	return false
}
