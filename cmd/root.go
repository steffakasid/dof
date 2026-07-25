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
var version string
var profileName string

var rootCmd = &cobra.Command{
	Use:   "dof",
	Short: "dof - dotfile repository manager",
	Long: `dof manages dotfiles using a bare git repository.

It uses a configuration file, environment variables, or command flags to select the repository and branch.

This tool requires git to be installed and available in PATH.`,
	TraverseChildren: true,
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
func Execute(v string) {
	version = v
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.dof.yaml)")
	rootCmd.PersistentFlags().StringVar(&profileName, "profile", "default", "Named profile to use from the config file")
	rootCmd.PersistentFlags().StringP("repository", "r", "", "Path to the bare git repository")
	rootCmd.PersistentFlags().StringP("branch", "b", "", "Git branch to track in the bare repository")

	cobra.OnInitialize(initConfig, initFlags)
}

func initFlags() {
	viper.SetDefault("repository", path.Join(userHomeDir, ".dof"))
	viper.SetDefault("branch", "main")
	viper.SetDefault("skip_files", []string{})

	// If a profile is selected, override defaults from profile config
	if profileName == "" {
		profileName = viper.GetString("default_profile")
	}
	if profileName != "" {
		logger.Infof("Using profile: %s", profileName)
		applyProfile(profileName)
	}

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

	logger.Debugf("repository: %s", viper.GetString("repository"))
	logger.Debugf("branch: %s", viper.GetString("branch"))
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
	// Merge global skip_files with profile-specific skip_files
	globalSkip := viper.GetStringSlice("skip_files")
	profileSkip := viper.GetStringSlice(profileKey + ".skip_files")
	if len(profileSkip) > 0 {
		merged := []string{}
		merged = append(merged, globalSkip...)
		merged = append(merged, profileSkip...)
		viper.Set("skip_files", deduplicate(merged))
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

func deduplicate(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
