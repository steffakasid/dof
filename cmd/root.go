// Package cmd implements the CLI commands for dof.
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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var version = "not set"
var profileName string

var rootCmd = &cobra.Command{
	Use:   "dof",
	Short: "dof - <do>t <f>ile repository tool",
	Long: `This tool is indended to setup and use a dot file repository. Basically the idea came
  when reading https://www.atlassian.com/git/tutorials/dotfiles. But finally I didn't like the
  way to do it with aliases (which must exist and be defined in a dotfile e.g. zshrc. Therefore
  to avoid this chicken and egg problem. I decided to write a little go program and here it is ;)

  The tool expects to use git from the path. So if you don't have git, it will not work!`,
	Version: version,
}

var (
	userHomeDir  string
	workDir      string
	repoPathName string
	gitAlias     *exec.Cmd
	logger       *Logger
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing config:", err)
		os.Exit(1)
	}
}

func init() {
	var err error
	if logger == nil {
		logger = NewOutputLogger(1)
	}
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting home directory:", err)
		os.Exit(1)
	}
	cobra.OnInitialize(initConfig, initFlags)
}

func initFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dof.yaml)")
	rootCmd.PersistentFlags().StringVar(&profileName, "profile", "", "Named profile to use (defined in config file)")

	viper.SetDefault("repository", path.Join(userHomeDir, ".dof"))
	viper.SetDefault("branch", "main")

	// If a profile is selected, override defaults from profile config
	if profileName != "" {
		applyProfile(profileName)
	}

	rootCmd.PersistentFlags().StringP("repository", "r", viper.GetString("repository"), "Repository folder to create a bare repository inside")
	if err := viper.BindPFlag("repository", rootCmd.PersistentFlags().Lookup("repository")); err != nil {
		fmt.Fprintln(os.Stderr, "Error binding repository flag:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(viper.GetString("repository"), 0700); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating repository directory:", err)
		os.Exit(1)
	}

	workDir, repoPathName = filepath.Split(viper.GetString("repository"))
	gitAlias = exec.Command("git", "--git-dir="+viper.GetString("repository"), "--work-tree="+workDir)

	logger.Debugln("branch:", viper.GetString("branch"))
	rootCmd.PersistentFlags().StringP("branch", "b", viper.GetString("branch"), "Set the branch to use")
	if err := viper.BindPFlag("branch", rootCmd.PersistentFlags().Lookup("branch")); err != nil {
		fmt.Fprintln(os.Stderr, "Error binding branch flag:", err)
		os.Exit(1)
	}
}

func applyProfile(name string) {
	profileKey := "profiles." + name
	if !viper.IsSet(profileKey) {
		fmt.Fprintf(os.Stderr, "Profile %q not found in config\n", name)
		os.Exit(1)
	}

	if repo := viper.GetString(profileKey + ".repository"); repo != "" {
		viper.Set("repository", repo)
	}
	if branch := viper.GetString(profileKey + ".branch"); branch != "" {
		viper.Set("branch", branch)
	}
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
