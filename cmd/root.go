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
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

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
	branch       string
	gitAlias     *exec.Cmd
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err := viper.SafeWriteConfig()
	if err != nil {
		log.Print(err)
	}
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)
	userHomeDir, err = os.UserHomeDir()
	doWePanic(err)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dof.yaml)")

	rootCmd.PersistentFlags().StringVarP(&repoPath, "repository", "r", path.Join(userHomeDir, ".dof"), "Repository folder to create a bare repository inside (default is $HOME/.dof)")
	viper.BindPFlag("repository", rootCmd.Flags().Lookup("repository"))
	err = os.MkdirAll(repoPath, 0700)
	doWePanic(err)
	workDir, repoPathName = filepath.Split(repoPath)
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+workDir)

	checkoutCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Set the branch to use (default is main)")
	viper.BindPFlag("branch", checkoutCmd.Flags().Lookup("branch"))

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
		log.Print("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Print(err)
		err := viper.SafeWriteConfigAs(viper.GetViper().ConfigFileUsed())
		if err != nil {
			log.Print(err)
		}
	}
}
