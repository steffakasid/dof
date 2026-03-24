# Requirements – dof (dotfile repository tool)

## Overview

**dof** is a CLI tool that manages dotfiles using a Git bare repository.
It wraps common git operations so users can version-control their
configuration files (e.g. `.zshrc`, `.gitconfig`) directly from their
home directory without turning `$HOME` into a regular git repo.

---

## Existing Features (implemented)

### REQ-01 Initialize a dotfile repository

> As a user, I want to initialize a new bare git repository so that I
> can start tracking my dotfiles.

**Acceptance criteria:**

- `dof init` creates a bare repo at the configured path (default `$HOME/.dof`).
- The configured branch is checked out.
- Untracked files are hidden (`status.showUntrackedFiles no`).
- A `.gitignore` entry for the repo directory is created and committed.

### REQ-02 Checkout an existing dotfile repository

> As a user, I want to clone an existing remote dotfile repository so
> that I can restore my dotfiles on a new machine.

**Acceptance criteria:**

- `dof checkout <git-repo-url>` clones the repo as a bare repo.
- Existing files that conflict are renamed with a `_before_dof` suffix.
- The configured branch is checked out.
- Untracked files are hidden.

### REQ-03 Add files to the repository

> As a user, I want to add a dotfile and commit it in one step.

**Acceptance criteria:**

- `dof add <file>` stages the file and commits it with message `Add <file>`.

### REQ-04 Synchronize with remote

> As a user, I want to push and pull changes to/from the remote in one
> command.

**Acceptance criteria:**

- `dof sync` commits all tracked changes, pushes, then pulls with rebase.
- `--push-only` / `-P` skips the pull step.
- `--pull-only` / `-p` skips the push step.
- If there are no local changes the push step is skipped gracefully.

### REQ-05 Show repository status

> As a user, I want a quick view of changed files.

**Acceptance criteria:**

- `dof status` prints short git status (`-s` flag) of the bare repo.

### REQ-06 Run arbitrary git commands

> As a user, I want to execute any git sub-command against my dotfile
> repository.

**Acceptance criteria:**

- `dof alias <git-args…>` forwards arguments to git with the correct
  `--git-dir` and `--work-tree` flags.

### REQ-07 Shell completion

> As a user, I want shell completion scripts for bash, zsh, fish, and
> powershell.

**Acceptance criteria:**

- `dof completion [bash|zsh|fish|powershell]` prints the completion
  script to stdout.

### REQ-08 Version information

> As a user, I want to see the current version and build time.

**Acceptance criteria:**

- `dof version` prints the version string and build timestamp.

### REQ-09 Configuration

> As a user, I want to configure the repository path and branch via
> flags, config file, or environment variables.

**Acceptance criteria:**

- `--repository` / `-r` overrides the repo path (default `$HOME/.dof`).
- `--branch` / `-b` overrides the branch name (default `main`).
- `--config` overrides the config file path (default `$HOME/.dof.yaml`).
- Environment variables prefixed with `DOF_` are respected.
- Config is persisted to `$HOME/.dof.yaml` on every run.

---

## Planned Features (not yet implemented)

### REQ-10 Set remote during init

> As a user, I want to provide a remote URL when running `dof init` so
> that I don't have to manually add it afterwards.

**Acceptance criteria:**

- `dof init --remote <url>` (or positional arg) adds the origin remote
  and sets upstream tracking after initializing the bare repo.

### REQ-11 Manage multiple dotfile repositories

> As a user, I want to manage more than one dotfile repository (e.g.
> work vs. personal) from one `dof` installation.

**Acceptance criteria:**

- A concept for named profiles or repos is introduced.
- Each profile has its own repo path, branch, and remote.
- The user can switch between profiles easily.
- Config file stores multiple profile entries.

### REQ-12 Refactor flags

> As a developer, I want the CLI flags to be clean and consistent so
> that maintenance is simpler.

**Acceptance criteria:**

- Global flags (`--repository`, `--branch`, `--config`) are reviewed
  and potentially reorganised.
- Conflicting or redundant flags are removed.
- Flag names follow a consistent naming convention.
- Verify that all defined flags work as intended
