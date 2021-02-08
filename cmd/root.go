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
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var version = "not set"

var rootCmd = &cobra.Command{
	Use:   "dof",
	Short: "dof - <do>t <f>ile repository tool",
	Long: `This tool is indended to setup and use a dot file repository. Basically the idea came
  when reading https://www.atlassian.com/git/tutorials/dotfiles. But finally I didn't like the
  way to do it with aliases (which must exist and be defined in a dotfile e.g. zshrc. Therefore
  to avoid this chicken and egg problem. I decided to write a little go program and here it is ;)

  The tool expects to use git from the path. So if you don't have git, it will not work!`,
}

var (
	userHomeDir  string
	repoPath     string
	workDir      string
	repoPathName string
	gitAlias     *exec.Cmd
	logger       *Logger
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}

	if err := viper.WriteConfig(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	var err error
	if logger == nil {
		logger = NewOutputLogger(1)
	}
	userHomeDir, err = os.UserHomeDir()
	doWePanic(err)
	cobra.OnInitialize(initConfig, initFlags)
}

func initFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dof.yaml)")

	viper.SetDefault("repository", path.Join(userHomeDir, ".dof"))
	rootCmd.PersistentFlags().StringP("repository", "r", viper.GetString("repository"), "Repository folder to create a bare repository inside")
	doWePanic(viper.BindPFlag("repository", rootCmd.PersistentFlags().Lookup("repository")))
	doWePanic(os.MkdirAll(viper.GetString("repository"), 0700))

	workDir, repoPathName = filepath.Split(viper.GetString("repository"))
	gitAlias = exec.Command("git", "--git-dir="+viper.GetString("repository"), "--work-tree="+workDir)

	viper.SetDefault("branch", "main")
	logger.Debugln("branch:", viper.GetString("branch"))
	rootCmd.PersistentFlags().StringP("branch", "b", viper.GetString("branch"), "Set the branch to use")
	doWePanic(viper.BindPFlag("branch", rootCmd.PersistentFlags().Lookup("branch")))
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
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		logger.Error(err)
		err := viper.SafeWriteConfig()
		if err != nil {
			logger.Fatal(err)
		}
	}
}
