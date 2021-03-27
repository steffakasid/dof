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
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/steffakasid/dof/internal"
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

  Example Usage:
  dot init
  dot alias remote add origin <git-repo-url>
  dot add .zshrc
  dot sync --push-only`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Initialize git bare repository...")
		// git init --bare $HOME/.cfg
		repo, err := git.PlainInit(repoPath, true)
		eh.IsFatalError(err)

		// TODO: do we need to set the remote directly here?
		logger.Infof("Checkout %s branch\n", viper.GetString("branch"))
		branch := &config.Branch{
			Name:   viper.GetString("branch"),
			Remote: "origin",
			Rebase: "true",
		}
		err = repo.CreateBranch(branch)
		eh.IsFatalError(err)

		doNotShowUntrackedFiles(repo)

		addGitIgnore()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	if logger == nil {
		logger = internal.NewOutputLogger(1)
	}
}

func addGitIgnore() {
	gitIgnore := path.Join(workDir, ".gitignore")
	file, err := os.Create(gitIgnore)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	linesToWrite := []string{repoPathName}
	for _, line := range linesToWrite {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	writer.Flush()

	addAndCommit(gitIgnore)
}

func doNotShowUntrackedFiles(repo *git.Repository) {
	// alias config='/usr/bin/git --git-dir=$HOME/.cfg/ --work-tree=$HOME'
	// config config --local status.showUntrackedFiles no
	cfg, err := repo.Config()
	eh.IsFatalError(err)
	cfg.Raw.SetOption("status", "", "showuntrackedfiles", "no")
	err = repo.SetConfig(cfg)
	eh.IsFatalError(err)
}
