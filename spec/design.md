# Design – dof (dotfile repository tool)

## 1. Architecture Overview

```text
┌────────────┐
│   main.go  │  entry point – calls cmd.Execute()
└─────┬──────┘
      │
      ▼
┌────────────────────────────────────────────────┐
│                   cmd/ package                 │
│                                                │
│  root.go      – cobra root command, viper cfg  │
│  init.go      – dof init                       │
│  checkout.go  – dof checkout                   │
│  add.go       – dof add                        │
│  sync.go      – dof sync                       │
│  status.go    – dof status                     │
│  alias.go     – dof alias (pass-through)       │
│  completion.go– shell completions              │
│  version.go   – version info                   │
│  release.go   – build-tag gated build metadata │
│  exec.go      – git command execution helpers  │
│  logger.go    – custom logrus wrapper          │
└────────────────────────────────────────────────┘
        │
        │  os/exec
        ▼
   ┌──────────┐
   │   git    │  system git binary
   └──────────┘
```

All commands shell out to the locally installed `git` binary via
`os/exec`. The bare repository lives at a configurable path (default
`$HOME/.dof`) and the work tree is the parent directory of that path
(typically `$HOME`).

## 2. Project Structure

```text
dof/
├── main.go                  # entry point
├── go.mod / go.sum
├── cmd/
│   ├── root.go              # cobra root, viper init
│   ├── init.go              # init command
│   ├── checkout.go          # checkout command
│   ├── add.go               # add command
│   ├── sync.go              # sync command
│   ├── status.go            # status command
│   ├── alias.go             # alias (pass-through) command
│   ├── completion.go        # shell completion
│   ├── version.go           # version command
│   ├── release.go           # build time (release tag only)
│   ├── exec.go              # helpers: execCmdAndPrint, execCmdAndReturn, doWePanic
│   └── logger.go            # Logger struct wrapping logrus
├── .goreleaser.yaml         # GoReleaser config
├── .github/
│   └── workflows/
│       ├── go-test.yml      # CI – go test
│       ├── release.yml      # semantic-release + goreleaser
│       └── codeql-analysis.yml
├── renovate.json            # Renovate dependency updates
├── spec/                    # spec-driven docs
│   ├── requirements.md
│   ├── design.md
│   └── tasks.md
└── README.adoc
```

## 3. Technology Stack

| Concern              | Choice                            | Rationale                                          |
|----------------------|-----------------------------------|----------------------------------------------------|
| Language             | Go 1.26                           | Existing project language                          |
| CLI framework        | spf13/cobra                       | De-facto standard for Go CLIs                      |
| Configuration        | spf13/viper                       | Seamless flag / env / file config binding          |
| Logging              | sirupsen/logrus (custom wrapper)  | Already in use; provides structured logging        |
| Git interaction      | os/exec → system `git`            | No cgo; leverages user's existing git installation |
| Build / Release      | GoReleaser + go-semantic-release  | Automated versioning, cross-compile, Homebrew tap  |
| CI                   | GitHub Actions                    | go-test, CodeQL, release workflows                 |
| Dependency updates   | Renovate                          | Automated PRs for Go modules & GH Actions          |

## 4. Configuration Management

Configuration is managed by **viper** with the following precedence
(highest → lowest):

1. CLI flags (`--repository`, `--branch`, `--config`)
2. Environment variables (`DOF_REPOSITORY`, `DOF_BRANCH`)
3. Config file (`$HOME/.dof.yaml`)
4. Defaults (`repository=$HOME/.dof`, `branch=main`)

On every successful run the merged config is written back to the YAML
file via `viper.WriteConfig()`.

## 5. Git Command Execution

A global `*exec.Cmd` template (`gitAlias`) is built once during flag
initialisation:

```text
git --git-dir=<repo-path> --work-tree=<work-dir>
```

Each command copies this template, appends its specific arguments, and
calls one of:

- `execCmdAndPrint(cmd)` – runs the command, prints stdout/stderr, panics on error.
- `execCmdAndReturn(cmd)` – runs the command and returns stdout as a string.

## 6. Error Handling Strategy

Currently a simple `doWePanic(err)` helper calls `logger.Fatal(err)`
(which calls `os.Exit(1)`) on any non-nil error. This is acceptable
for a CLI tool but could be improved by returning errors up the call
chain and letting cobra handle exit codes.

## 7. Security Considerations

- No secrets are stored in code; tokens (`GITHUB_TOKEN`,
  `HOMEBREW_TAP_GITHUB_TOKEN`) are injected via GitHub Actions secrets.
- The repo directory is created with mode `0700`.
- The tool passes user-supplied arguments to `git` via `os/exec`
  (not through a shell), which avoids shell injection.
- CodeQL analysis runs weekly and on every PR.

## 8. Build & Release

- **GoReleaser** builds for linux, windows, darwin (all with
  `CGO_ENABLED=0`).
- The `release` build tag injects `buildTime` via `cmd/release.go`.
- Version is injected via ldflags: `-X cmd.version={{.Version}}`.
- A Homebrew formula is pushed to `steffakasid/homebrew-dof`.
- Release is triggered on a monthly schedule via `go-semantic-release`.
