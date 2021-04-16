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
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/steffakasid/dof/internal"

	"github.com/spf13/viper"
)

var cfgFile string
var version = "not set"

var rootCmd = &cobra.Command{
	Use:     "dof",
	Short:   "dof - <do>t <f>ile repository tool",
	Version: version,
	Long: `This tool is indended to setup and use a dot file repository. Basically the idea came
when reading https://www.atlassian.com/git/tutorials/dotfiles. But finally I didn't like the
way to do it with aliases (which must exist and be defined in a dotfile e.g. zshrc. Therefore
to avoid this chicken and egg problem. I decided to write a little go program and here it is ;)

Most of the git commands now uses go-git library but as the status command doesn't work like
the native one a local git installation which is in the path is still necessary.

If you experience any issues which can't be resolved via dof you can still run any git command
using:
/usr/bin/git --git-dir=$HOME/.dof/ --work-tree=$HOME`,
}

var (
	userHomeDir    string
	workDir        string
	repoFolderName string
	logger         *internal.Logger
	traceLogger    *internal.Logger
	eh             *internal.ErrorHandler
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
	if viper.ConfigFileUsed() != "" {
		if err := viper.WriteConfig(); err != nil {
			logger.Fatal(err)
		}
	}
}

func init() {
	var err error
	viper.SetDefault("LogLevel", 4)
	traceLogger = internal.NewTraceLogger(logrus.Level(viper.GetInt("LogLevel")), 2)

	if logger == nil {
		logger = internal.NewOutputLogger(1)
	}
	if eh == nil {
		eh = internal.NewErrorHandler(logger)
	}
	userHomeDir, err = os.UserHomeDir()
	eh.IsFatalError(err)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dof.yaml)")

	viper.SetDefault("repository", path.Join(userHomeDir, ".dof"))
	rootCmd.PersistentFlags().StringP("repository", "r", viper.GetString("repository"), "Repository folder to create a bare repository inside")
	eh.IsFatalError(viper.BindPFlag("repository", rootCmd.PersistentFlags().Lookup("repository")))
	eh.IsFatalError(os.MkdirAll(viper.GetString("repository"), 0700))

	workDir, repoFolderName = filepath.Split(viper.GetString("repository"))

	viper.SetDefault("branch", "main")
	traceLogger.Debugln("branch:", viper.GetString("branch"))
	rootCmd.PersistentFlags().StringP("branch", "b", viper.GetString("branch"), "Set the branch to use")
	eh.IsFatalError(viper.BindPFlag("branch", rootCmd.PersistentFlags().Lookup("branch")))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(userHomeDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dof")
	}
	viper.SetEnvPrefix("DOF_")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		traceLogger.SetLevel(logrus.Level(viper.GetInt("LogLevel")))

		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		eh.IsError(err)
		err := viper.SafeWriteConfig()
		if err != nil {
			eh.IsFatalError(err)
		}
	}
}
