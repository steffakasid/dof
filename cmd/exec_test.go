package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if logger == nil {
		logger = NewOutputLogger(1)
	}
}

func TestBuildGitEnv(t *testing.T) {
	t.Run("empty git_env returns os.Environ unchanged", func(t *testing.T) {
		viper.Reset()
		viper.SetDefault("git_env", map[string]string{})

		env := buildGitEnv()
		assert.ElementsMatch(t, os.Environ(), env)
	})

	t.Run("git_env entries are appended on top of os.Environ", func(t *testing.T) {
		viper.Reset()
		viper.Set("git_env", map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"})

		env := buildGitEnv()
		assert.Contains(t, env, "GIT_CONFIG_GLOBAL=~/.gitconfig")
		// existing OS environment is retained
		for _, e := range os.Environ() {
			assert.Contains(t, env, e)
		}
	})

	t.Run("values are passed verbatim without expansion", func(t *testing.T) {
		viper.Reset()
		viper.Set("git_env", map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"})

		env := buildGitEnv()
		assert.Contains(t, env, "GIT_CONFIG_GLOBAL=~/.gitconfig")
		// ensure "~" was not expanded to a home directory path
		for _, e := range env {
			assert.NotContains(t, e, "GIT_CONFIG_GLOBAL="+os.Getenv("HOME"))
		}
	})
}

func TestNewGitCmd(t *testing.T) {
	viper.Reset()
	viper.Set("git_env", map[string]string{"GIT_CONFIG_GLOBAL": "~/.gitconfig"})

	cmd := newGitCmd("status")
	require.NotNil(t, cmd)
	assert.Equal(t, "git", cmd.Args[0])
	assert.Equal(t, "status", cmd.Args[1])
	assert.Contains(t, cmd.Env, "GIT_CONFIG_GLOBAL=~/.gitconfig")
}

func TestExecCmdAndPrint(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *exec.Cmd
		wantErr bool
	}{
		{
			name:    "successful command",
			cmd:     exec.Command("echo", "hello"),
			wantErr: false,
		},
		{
			name:    "failing command",
			cmd:     exec.Command("false"),
			wantErr: true,
		},
		{
			name:    "command not found",
			cmd:     exec.Command("nonexistent-command-xyz"),
			wantErr: true,
		},
		{
			name:    "command with stderr output",
			cmd:     exec.Command("sh", "-c", "echo error >&2 && exit 0"),
			wantErr: false,
		},
		{
			name:    "command with stderr and failure",
			cmd:     exec.Command("sh", "-c", "echo error >&2 && exit 1"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := execCmdAndPrint(tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExecCmdAndReturn(t *testing.T) {
	tests := []struct {
		name       string
		cmd        *exec.Cmd
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "successful command returns output",
			cmd:        exec.Command("echo", "hello"),
			wantOutput: "hello\n",
			wantErr:    false,
		},
		{
			name:       "failing command returns error",
			cmd:        exec.Command("false"),
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:       "command not found returns error",
			cmd:        exec.Command("nonexistent-command-xyz"),
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:       "command with multi-line output",
			cmd:        exec.Command("printf", "line1\nline2\n"),
			wantOutput: "line1\nline2\n",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := execCmdAndReturn(tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, output)
				assert.Contains(t, err.Error(), "failed")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOutput, output)
			}
		})
	}
}
