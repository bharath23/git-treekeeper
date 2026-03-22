![](https://github.com/bharath23/git-treekeeper/actions/workflows/test-suite.yml/badge.svg)

# Git TreeKeeper

**Git TreeKeeper** is a CLI tool to help manage Git branches and worktrees
more efficiently.

## Table of Contents

- [Overview](#overview)
- [CLI](#cli)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Shell Integration](#shell-integration)
- [Examples](#examples)
- [Build and Install](#build-and-install)
- [Install from Release](#install-from-release)
- [License](#license)

## Overview

Git TreeKeeper focuses on:

- Cloning repositories
- Creating new branches and worktrees
- Checking out branches
- Deleting branches and worktrees
- Listing worktrees
- Checking worktree health (including stale and orphaned detection)
- Pruning stale worktrees
- Garbage collecting old merged branches (dry-run by default)
- Syncing branches with remotes (including `sync --all`)
- Workflow setup for forks
- Shell autocompletion for Bash/Zsh

Upcoming Roadmap:

- **Merge/Rebase Management**: Proactive tools to resolve stuck merges or rebases.
- **Auto-Fix**: `git tk doctor --fix` to automatically resolve worktree health issues.
- **Stacked Branches**: Support for managing chains of dependent PRs.
- **Ephemeral Worktrees**: Short-lived, self-cleaning workspaces for CI/CD and reviews.

## CLI

The binary is named `git-tk`. Git discovers it as a plugin, so users invoke it
as `git tk ...`. Direct invocation with `git-tk ...` also works.

**Pass-Through:** Any command or flag not explicitly handled by `git-tk` is
automatically passed through to the underlying `git` command. For example,
`git tk status` runs `git status`, and `git tk branch -v` runs `git branch -v`.

## Quick Start

```bash
git tk clone <repo-url>
cd repo/worktrees/main
git tk branch <branch-name>
git tk checkout <branch-name>
git tk list
git tk doctor
git tk prune
git tk sync --default
git tk setup --upstream-url <url>
```

## Command Reference

Full command behavior and flags are documented in [docs/commands.md](docs/commands.md).

Common commands:
- `git tk branch`
- `git tk checkout`
- `git tk clone`
- `git tk repair`
- `git tk doctor`
- `git tk list`
- `git tk prune`
- `git tk gc`
- `git tk sync`
- `git tk setup`

## Shell Integration

Use `--path-only` to auto-`cd` into the returned worktree path. Example for bash/zsh:

```bash
gtk() {
  local gittk
  gittk="$(command -v git-tk)" || { echo "git-tk not found in PATH"; return 127; }
  if [ "$1" = "branch" ] || [ "$1" = "checkout" ] || [ "$1" = "clone" ]; then
    local wt_path
    wt_path="$($gittk "$@" --path-only)"
    [ -d "$wt_path" ] && cd "$wt_path"
  else
    "$gittk" "$@"
  fi
}
```

## Examples

For detailed, scenario-based examples, see [docs/examples.md](docs/examples.md).

## Build and Install

```bash
make build
```

The binary will be written to `build/git-tk` by default. You can override:

- `BUILD_DIR` (output directory)
- `BIN` (binary name)
- `VERSION` (ldflags version)

```bash
make install
```

`make install` installs to the default Go install location. Override with
`GOBIN` (preferred) or `GOPATH/bin`, and you can also override `VERSION`.

Examples:

```bash
make build BUILD_DIR=dist BIN=git-tk
make build VERSION=v0.1.0
GOBIN="$HOME/.local/bin" make install
```

To run tests:

```bash
make test
```

For verbose local output:

```bash
make test V=1
```

Short alias:

```bash
make test-v
```

CI runs:

```bash
make test-ci
```

To run formatting and vet checks:

```bash
make check
```

## Install from Release

1. Download the correct asset for your OS/arch from GitHub Releases:
   `git-tk-<version>-<os>-<arch>` and its matching `.sha256`.
2. Verify the checksum:

```bash
# macOS
shasum -a 256 git-tk-<version>-<os>-<arch>
cat git-tk-<version>-<os>-<arch>.sha256
```

```bash
# Linux (or macOS with coreutils installed)
sha256sum -c git-tk-<version>-<os>-<arch>.sha256
```

3. Make it executable and move it into your `PATH`:

```bash
chmod +x git-tk-<version>-<os>-<arch>
mv git-tk-<version>-<os>-<arch> /usr/local/bin/git-tk
```

On macOS you may also need to allow execution in System Settings > Privacy & Security
on first run.

## License

MIT License. See [LICENSE](LICENSE) file.
