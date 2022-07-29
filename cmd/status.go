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
	"github.com/steffakasid/dof/internal"
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
	logger.Info("Status of dof repository:")
	traceLogger.Debug("repoPath: ", repoFolderName)
	traceLogger.Debug("workDir:", workDir)

	dofRepo, err := internal.OpenDofRepo(workDir, repoFolderName)
	eh.IsFatalError(err)

	out, err := dofRepo.Status()
	eh.IsFatalError(err)
	if len(out) > 0 {
		logger.Info(string(out))
	} else {
		logger.Info("No files changed!")
	}

}

func init() {
	rootCmd.AddCommand(statusCmd)
	if logger == nil {
		logger = internal.NewOutputLogger(1)
	}
}
