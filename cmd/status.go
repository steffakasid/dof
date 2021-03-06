package cmd

/*
Copyright © 2021 Steffen Rumpf <github@steffen-rumpf.de>

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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Args:  cobra.NoArgs,
	Short: "Get git status of repository",
	Long:  `Get git status of repository`,
	Run:   status,
}

func status(cmd *cobra.Command, args []string) {
	gitStatus := *gitAlias
	gitStatusArgs := []string{"status", "-s"}
	gitStatus.Args = append(gitStatus.Args, gitStatusArgs...)
	logger.Info("Status of dof repository...")
	execCmdAndPrint(&gitStatus)
}

func init() {
	rootCmd.AddCommand(statusCmd)
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}
