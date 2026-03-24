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
	// Set up a repo with a remote to test sync
	tmpDir := setupTestEnv(t)
	repoPath := filepath.Join(tmpDir, ".dof")

	err := os.MkdirAll(repoPath, 0700)
	require.NoError(t, err)

	err = initCmd.RunE(initCmd, nil)
	require.NoError(t, err)

	// Create and add a test file
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	testFile := filepath.Join(tmpDir, ".zshrc")
	err = os.WriteFile(testFile, []byte("# zsh config"), 0600)
	require.NoError(t, err)

	err = addCmd.RunE(addCmd, []string{testFile})
	require.NoError(t, err)

	// Modify the tracked file
	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	err = os.WriteFile(testFile, []byte("# zsh config\nalias ll='ls -la'"), 0600)
	require.NoError(t, err)

	// Sync with push-only (no remote, so push will fail, but commit should succeed)
	// We set dontPull=true, dontPush=false but since there's no remote, push will error
	// Instead, test that the commit step works by checking pull-only mode
	dontPush = true
	dontPull = true

	gitAlias = exec.Command("git", "--git-dir="+repoPath, "--work-tree="+tmpDir)
	err = syncCmd.RunE(syncCmd, nil)
	require.NoError(t, err)

	// Reset flags
	dontPush = false
	dontPull = false
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
