package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnv configures package-level variables for a test
// using a temporary directory as the working tree.
func setupTestEnv(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Save and restore cwd since some commands call os.Chdir
	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	repoPath := filepath.Join(tmpDir, ".dof")
	require.NoError(t, os.MkdirAll(repoPath, 0700))

	viper.Set("repository", repoPath)
	viper.Set("branch", "main")

	workDir = tmpDir + string(os.PathSeparator)
	repoPathName = ".dof"
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)

	logger = NewOutputLogger(1)

	return tmpDir
}

func TestInitCommand(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Verify bare repo was created
	headFile := filepath.Join(repoPath, "HEAD")
	assert.FileExists(t, headFile)

	// Verify .gitignore was created and committed
	gitIgnore := filepath.Join(tmpDir, ".gitignore")
	assert.FileExists(t, gitIgnore)

	content, err := os.ReadFile(gitIgnore)
	require.NoError(t, err)
	assert.Contains(t, string(content), ".dof")
}

func TestAddCommand(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	// Initialize repo first
	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Create a test file to add
	testFile := filepath.Join(tmpDir, ".testrc")
	err = os.WriteFile(testFile, []byte("test config"), 0600)
	require.NoError(t, err)

	// Reset gitAlias since init modifies it
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)

	err = addCmd.RunE(addCmd, []string{testFile})
	require.NoError(t, err)

	// Verify the file was committed by checking git log
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	gitLog := *gitAlias
	gitLog.Args = append(gitLog.Args, "log", "--oneline")
	output, err := execCmdAndReturn(&gitLog)
	require.NoError(t, err)
	assert.Contains(t, output, "Add")
}

func TestStatusCommand(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	// Initialize repo first
	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Reset gitAlias
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)

	// Status should work without error on a clean repo
	err = statusCmd.RunE(statusCmd, nil)
	require.NoError(t, err)
}

func TestCheckoutCommand(t *testing.T) {
	// Save and restore cwd since checkout calls os.Chdir
	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	// Phase 1: Create a source bare repo with some content
	srcDir := t.TempDir()
	srcRepo := filepath.Join(srcDir, ".dof")

	// Set up env for source init
	viper.Set("repository", srcRepo)
	viper.Set("branch", "main")
	workDir = srcDir + string(os.PathSeparator)
	repoPathName = ".dof"
	gitAlias = exec.Command("git", "--git-dir="+srcRepo, "--work-tree="+srcDir)
	logger = NewOutputLogger(1)

	err = os.MkdirAll(srcRepo, 0700)
	require.NoError(t, err)

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Phase 2: Checkout from source into a new temp dir
	destDir := t.TempDir()
	destRepo := filepath.Join(destDir, ".dof")

	viper.Set("repository", destRepo)
	workDir = destDir + string(os.PathSeparator)
	repoPathName = ".dof"
	gitAlias = exec.Command("git", "--git-dir="+destRepo, "--work-tree="+destDir)

	err = checkoutCmd.RunE(checkoutCmd, []string{srcRepo})
	require.NoError(t, err)

	// Verify the bare repo was cloned
	headFile := filepath.Join(destRepo, "HEAD")
	assert.FileExists(t, headFile)
}

func TestSyncCommandPushOnly(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	// Create a bare remote repo to push to
	remoteDir := t.TempDir()
	remoteRepo := filepath.Join(remoteDir, "remote.git")
	gitInitBare := exec.Command("git", "init", "--bare", remoteRepo)
	require.NoError(t, gitInitBare.Run())

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Add remote origin
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	gitRemoteAdd := *gitAlias
	gitRemoteAdd.Args = append(gitRemoteAdd.Args, "remote", "add", "origin", remoteRepo)
	require.NoError(t, gitRemoteAdd.Run())

	// Create and add a test file
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	testFile := filepath.Join(tmpDir, ".zshrc")
	err = os.WriteFile(testFile, []byte("# zsh config"), 0600)
	require.NoError(t, err)

	err = addCmd.RunE(addCmd, []string{testFile})
	require.NoError(t, err)

	// Modify the tracked file so sync has something to commit and push
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	err = os.WriteFile(testFile, []byte("# zsh config\nalias ll='ls -la'"), 0600)
	require.NoError(t, err)

	// Test push-only mode: should commit and push, but not pull
	pushOnly = true
	pullOnly = false

	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	err = syncCmd.RunE(syncCmd, nil)
	require.NoError(t, err)

	// Reset flags
	pushOnly = false
	pullOnly = false
}

func TestInitCommandWithRemote(t *testing.T) {
	// Phase 1: Create a source bare repo to use as remote
	srcDir := t.TempDir()
	srcRepo := filepath.Join(srcDir, "remote.git")

	gitInitBare := exec.Command("git", "init", "--bare", srcRepo)
	err := gitInitBare.Run()
	require.NoError(t, err)

	// Phase 2: Init with --remote pointing to srcRepo
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err = os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	remoteURL = srcRepo
	t.Cleanup(func() { remoteURL = "" })

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Verify remote was added
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	gitRemoteV := *gitAlias
	gitRemoteV.Args = append(gitRemoteV.Args, "remote", "-v")
	output, err := execCmdAndReturn(&gitRemoteV)
	require.NoError(t, err)
	assert.Contains(t, output, "origin")
	assert.Contains(t, output, srcRepo)
}

func TestApplyProfile(t *testing.T) {
	tests := []struct {
		name        string
		profile     string
		setupConfig func()
		wantRepo    string
		wantBranch  string
		wantErr     bool
	}{
		{
			name:    "profile overrides repository and branch",
			profile: "work",
			setupConfig: func() {
				viper.Set("profiles.work.repository", "/tmp/work-dof")
				viper.Set("profiles.work.branch", "develop")
			},
			wantRepo:   "/tmp/work-dof",
			wantBranch: "develop",
		},
		{
			name:    "profile overrides only repository",
			profile: "personal",
			setupConfig: func() {
				viper.Set("profiles.personal.repository", "/tmp/personal-dof")
				viper.SetDefault("branch", "main")
			},
			wantRepo:   "/tmp/personal-dof",
			wantBranch: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.SetDefault("repository", "/default/.dof")
			viper.SetDefault("branch", "main")
			viper.SetDefault("skip_files", []string{})
			tt.setupConfig()

			applyProfile(tt.profile)

			assert.Equal(t, tt.wantRepo, viper.GetString("repository"))
			assert.Equal(t, tt.wantBranch, viper.GetString("branch"))
		})
	}
}

func TestApplyProfileWithSkipFiles(t *testing.T) {
	viper.Reset()
	viper.SetDefault("repository", "/default/.dof")
	viper.SetDefault("branch", "main")
	viper.SetDefault("skip_files", []string{})
	viper.Set("profiles.work.repository", "/tmp/work-dof")
	viper.Set("profiles.work.skip_files", []string{"README.md", "LICENSE"})

	applyProfile("work")

	assert.Equal(t, []string{"README.md", "LICENSE"}, viper.GetStringSlice("skip_files"))
}

func TestApplyProfileGitEnvMerge(t *testing.T) {
	tests := []struct {
		name    string
		global  map[string]string
		profile map[string]string
		want    map[string]string
	}{
		{
			name:    "profile only",
			global:  map[string]string{},
			profile: map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig_work"},
			want:    map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig_work"},
		},
		{
			name:    "profile wins per key",
			global:  map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"},
			profile: map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig_work"},
			want:    map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig_work"},
		},
		{
			name:    "global-only key merged in",
			global:  map[string]string{"GIT_SSH_COMMAND": "ssh -i ~/.ssh/id"},
			profile: map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig_work"},
			want: map[string]string{
				"GIT_SSH_COMMAND":   "ssh -i ~/.ssh/id",
				"GIT_CONFIG_GLOBAL": "~/.gitconfig_work",
			},
		},
		{
			name:    "empty profile keeps global",
			global:  map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"},
			profile: map[string]string{},
			want:    map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.SetDefault("repository", "/default/.dof")
			viper.SetDefault("branch", "main")
			viper.SetDefault("skip_files", []string{})
			viper.SetDefault("git_env", map[string]string{})
			viper.Set("git_env", tt.global)
			viper.Set("profiles.work.repository", "/tmp/work-dof")
			if len(tt.profile) > 0 {
				viper.Set("profiles.work.git_env", tt.profile)
			}

			applyProfile("work")

			assert.Equal(t, tt.want, viper.GetStringMapString("git_env"))
		})
	}
}

func TestInitHonorsGitEnv(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	// A git config file injected via GIT_CONFIG_GLOBAL with a distinctive value.
	gitConfigFile := filepath.Join(tmpDir, "custom_gitconfig")
	err := os.WriteFile(gitConfigFile, []byte("[user]\n\tname = From Git Env\n\temail = env@example.com\n"), 0600)
	require.NoError(t, err)

	viper.Set("git_env", map[string]string{"GIT_CONFIG_GLOBAL": gitConfigFile})
	t.Cleanup(func() { viper.Set("git_env", map[string]string{}) })

	// Rebuild gitAlias with the git_env applied (mirrors initFlags).
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	gitAlias.Env = buildGitEnv()

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// The bare repo must exist, proving git init --bare ran with the env.
	assert.FileExists(t, filepath.Join(repoPath, "HEAD"))

	// Verify git actually read GIT_CONFIG_GLOBAL: user.name resolves to our value
	// (the local identity is only set as a fallback when none is configured).
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	gitAlias.Env = buildGitEnv()
	gitConfigGet := *gitAlias
	gitConfigGet.Args = append(gitConfigGet.Args, "config", "--global", "user.name")
	out, err := execCmdAndReturn(&gitConfigGet)
	require.NoError(t, err)
	assert.Contains(t, out, "From Git Env")
}

func TestCheckoutHonorsGitEnv(t *testing.T) {
	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	// Phase 1: create a source bare repo with content.
	srcDir := t.TempDir()
	srcRepo := filepath.Join(srcDir, ".dof")
	viper.Set("repository", srcRepo)
	viper.Set("branch", "main")
	viper.Set("git_env", map[string]string{})
	workDir = srcDir + string(os.PathSeparator)
	repoPathName = ".dof"
	gitAlias = exec.Command("git", "--git-dir="+srcRepo, "--work-tree="+srcDir)
	logger = NewOutputLogger(1)
	require.NoError(t, os.MkdirAll(srcRepo, 0700))
	require.NoError(t, initCmd.RunE(initCmd, nil))

	// Phase 2: checkout into a new dir with a git_env configured.
	destDir := t.TempDir()
	destRepo := filepath.Join(destDir, ".dof")

	gitConfigFile := filepath.Join(destDir, "custom_gitconfig")
	require.NoError(t, os.WriteFile(gitConfigFile, []byte("[init]\n\tdefaultBranch = main\n"), 0600))

	viper.Set("repository", destRepo)
	viper.Set("git_env", map[string]string{"GIT_CONFIG_GLOBAL": gitConfigFile})
	t.Cleanup(func() { viper.Set("git_env", map[string]string{}) })
	workDir = destDir + string(os.PathSeparator)
	repoPathName = ".dof"
	gitAlias = exec.Command("git", "--git-dir="+destRepo, "--work-tree="+destDir)
	gitAlias.Env = buildGitEnv()

	err = checkoutCmd.RunE(checkoutCmd, []string{srcRepo})
	require.NoError(t, err)

	// The cloned bare repo must exist, proving git clone --bare ran with the env.
	assert.FileExists(t, filepath.Join(destRepo, "HEAD"))
}

func TestIsSkipped(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		skipFiles []string
		want      bool
	}{
		{
			name:      "file is in skip list",
			file:      "README.md",
			skipFiles: []string{"README.md", "LICENSE"},
			want:      true,
		},
		{
			name:      "file is not in skip list",
			file:      ".zshrc",
			skipFiles: []string{"README.md", "LICENSE"},
			want:      false,
		},
		{
			name:      "empty skip list",
			file:      "README.md",
			skipFiles: []string{},
			want:      false,
		},
		{
			name:      "nil skip list",
			file:      "README.md",
			skipFiles: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isSkipped(tt.file, tt.skipFiles))
		})
	}
}

func TestApplySkipFilesEmpty(t *testing.T) {
	setupTestEnv(t)
	viper.Set("skip_files", []string{})

	err := applySkipFiles()
	require.NoError(t, err)
}

func TestApplySkipFiles(t *testing.T) {
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Add a README.md to the repo
	readmeFile := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(readmeFile, []byte("# My Dotfiles"), 0600)
	require.NoError(t, err)

	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	err = addCmd.RunE(addCmd, []string{readmeFile})
	require.NoError(t, err)

	// Now configure skip_files and apply
	viper.Set("skip_files", []string{"README.md"})
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)

	err = applySkipFiles()
	require.NoError(t, err)

	// Verify README.md was removed from working tree
	assert.NoFileExists(t, readmeFile)

	// Verify sparse-checkout file was written
	sparseFile := filepath.Join(repoPath, "info", "sparse-checkout")
	assert.FileExists(t, sparseFile)

	content, err := os.ReadFile(sparseFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "/*")
	assert.Contains(t, string(content), "!README.md")
}
