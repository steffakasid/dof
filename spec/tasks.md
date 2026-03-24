# Implementation Plan – dof

Legend: `[x]` = done, `[ ]` = to do

---

## 1. Project Setup

- [x] **[Setup]** Initialise Go module and cobra/viper skeleton
  _Done when:_ `go build ./...` succeeds and `dof --help` prints usage.

- [x] **[Setup]** Add GoReleaser configuration
  _Done when:_ `.goreleaser.yaml` exists and `goreleaser check` passes.

- [x] **[Setup]** Add GitHub Actions workflows (test, release, CodeQL)
  _Done when:_ CI runs on push to main and on PRs.

- [x] **[Setup]** Add Renovate configuration
  _Done when:_ `renovate.json` exists with Go module and GH Action regex managers.

---

## 2. Core Commands (existing)

- [x] **[cmd/init]** Implement `dof init` – bare repo setup, branch checkout, .gitignore
  _Done when:_ Running `dof init` creates a bare repo, hides untracked files, and commits `.gitignore`.

- [x] **[cmd/checkout]** Implement `dof checkout` – clone, rename conflicts, checkout
  _Done when:_ Running `dof checkout <url>` clones the repo, renames conflicting files, and checks out the branch.

- [x] **[cmd/add]** Implement `dof add` – stage and commit a file
  _Done when:_ Running `dof add <file>` stages and commits the file.

- [x] **[cmd/sync]** Implement `dof sync` – commit, push, pull with `--push-only` / `--pull-only`
  _Done when:_ Running `dof sync` synchronises with remote; flags control direction.

- [x] **[cmd/status]** Implement `dof status` – short git status
  _Done when:_ Running `dof status` prints short status output.

- [x] **[cmd/alias]** Implement `dof alias` – pass-through git commands
  _Done when:_ Running `dof alias <args>` executes git with correct `--git-dir` / `--work-tree`.

- [x] **[cmd/completion]** Implement shell completion (bash, zsh, fish, powershell)
  _Done when:_ Running `dof completion bash` prints valid completion script.

- [x] **[cmd/version]** Implement version command
  _Done when:_ Running `dof version` prints version and build time.

---

## 3. Configuration

- [x] **[Config]** Viper-based config with flags, env vars, and YAML file
  _Done when:_ `--repository`, `--branch`, `--config` flags work; `DOF_` env vars are read; config persists to `$HOME/.dof.yaml`.

---

## 4. Testing

- [x] **[Testing]** Add unit tests for `execCmdAndPrint` and `execCmdAndReturn`
  _Done when:_ Tests exist in `cmd/exec_test.go`; they cover success and error cases using mocked commands.

- [ ] **[Testing]** Add unit tests for `Logger`
  _Done when:_ Tests exist in `cmd/logger_test.go`; they validate trace vs. output logger behaviour.

- [ ] **[Testing]** Add integration tests for core commands (init, checkout, add, sync, status)
  _Done when:_ Tests exist that exercise each command against a temporary git repo; coverage ≥ 80 % for core logic.

---

## 5. Planned Features

- [ ] **[cmd/init]** REQ-10 – Accept `--remote` flag during init
  _Done when:_ `dof init --remote <url>` adds origin and sets upstream; tested.

- [ ] **[Config]** REQ-11 – Multi-repo / profile support
  _Done when:_ Config file supports named profiles; `--profile` flag selects the active profile; each profile has its own repo path, branch, and remote; tested.

- [ ] **[Flags]** REQ-12 – Refactor global flags for consistency
  _Done when:_ Flags follow a consistent naming convention; no redundant flags remain; existing behaviour is preserved; all defined flags are verified to work as intended; tested.

---

## 6. Code Quality

- [x] **[Quality]** Replace `doWePanic` with proper error returns
  _Done when:_ Commands return errors via cobra's `RunE`; `doWePanic` is removed; no `log.Fatal` outside of `Execute()`.

- [x] **[Quality]** Add golangci-lint configuration and fix lint issues
  _Done when:_ `.golangci.yml` exists; `golangci-lint run` passes; CI workflow runs lint.
