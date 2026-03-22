![](https://github.com/bharath23/git-treekeeper/actions/workflows/test-suite.yml/badge.svg)

# Git TreeKeeper

**Git TreeKeeper** is a CLI tool to help manage Git branches and worktrees
efficiently.

## Overview

For now, Git TreeKeeper focuses on:

- Cloning repositories
- Creating new branches and worktrees
- Checking out branches
- Deleting branches and worktrees
- List worktrees
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
This ensures a seamless experience when using standard Git workflows.

Example usage:

```bash
git tk clone <repo-url>
git tk branch <branch-name>
git tk checkout <branch-name>
```

## Command Behavior

### branch

Creates a new branch and a worktree for it under `worktrees/<branch>`. If no
base is provided, it uses the current branch, and falls back to the default
branch when needed.

Deletion:
- `-d, --delete` deletes the branch and its worktree if it is merged.
- `-D, --force` deletes even if unmerged (still refuses dirty/in-progress).
- `--remote` also deletes the remote branch from `origin` (errors if missing).
- `--yes` skips confirmation for destructive deletes (`-D` or `--remote`).
- Deletion is refused if run inside the target branch worktree.

```bash
git tk branch <branch-name> [base]
```

Options:
- `--path-only` prints only the worktree path.

### checkout

Ensures a worktree exists for the branch and prints its path. If the branch
does not exist, it creates it from the default branch.

```bash
git tk checkout <branch-name>
```

Options:
- `--path-only` prints only the worktree path.

### clone

Clones a repository into the default directory (same as `git clone`), creates a
bare repo at `<base>/repo.git`, and adds a worktree for the default branch at
`<base>/worktrees/<branch>`.

```bash
git tk clone <repo-url> [path]
```

Options:
- `--path-only` prints only the worktree path.

### doctor

Reports worktree health. Detects:
- `clean` / `dirty` states.
- `merge` or `rebase` operations in progress.
- `stale (directory missing)`: worktrees tracked by Git but missing from disk.
- `orphaned directory`: directories in `worktrees/` that are not tracked by Git.
Tracking shows the upstream ref (for example `origin/main`) or `none`.

```bash
git tk doctor
```

Options:
- `--porcelain` prints tab-separated rows with no header.
- `--json` prints a JSON array of `{branch,state,tracking}` objects.

### list

Lists all worktrees with their branch name and path, sorted by branch name.

```bash
git tk list
```

Options:
- `--porcelain` prints tab-separated rows with no header.
- `--json` prints a JSON array of `{branch,path}` objects.

### prune

Prunes worktree entries whose directories no longer exist. By default it only
removes stale worktrees; use `--merged-branches` to also remove merged branches
that no longer have a worktree.

```bash
git tk prune
```

Options:
- `--dry-run` prints what would be pruned without making changes.
- `--merged-branches` also deletes merged branches that have no worktree.

### gc

Garbage collects local branches that are merged and older than the threshold.
Runs in dry-run mode by default and skips active worktrees.

```bash
git tk gc
git tk gc --apply
```

Options:
- `--apply` deletes branches (default is dry-run).
- `--age-days <days>` only considers branches older than the threshold (default 30).

### sync

Syncs a branch using `git fetch` plus `--ff-only` fast-forward, matching a
read-only main workflow. By default, it syncs only "clean" worktrees.

```bash
git tk sync
git tk sync --default
git tk sync --branch feature-x
git tk sync --all
```

Options:
- `--all` syncs all active, clean worktrees.
- `--default` syncs the default branch worktree.
- `--branch <name>` syncs the specified branch worktree.
- `--remote <name>` fetches from a specific remote.
- `--add-upstream <url>` adds an upstream remote (default name: `upstream`).
- `--set-upstream` sets the branch upstream to `<upstream>/<branch>` and sets
  push remote to `origin`.
- `--upstream <name>` sets the upstream remote name (default `upstream`).
- `--origin <name>` sets the origin remote name for push defaults.
- `--dry-run` prints what would be synced without making changes.

### setup

Configures fork workflow defaults by adding upstream (if URL provided),
setting the branch to track upstream, and setting push defaults to origin.

```bash
git tk setup --upstream-url <url>
git tk setup --branch main --upstream upstream --origin origin
```

Options:
- `--branch <name>` configures the specified branch (defaults to default branch).
- `--upstream <name>` sets the upstream remote name (default `upstream`).
- `--origin <name>` sets the origin remote name for push defaults.
- `--upstream-url <url>` adds the upstream remote if missing.
- `--install-hooks` installs a pre-commit hook blocking commits on the branch.
- `--dry-run` prints what would be configured without making changes.

### Global Flags

- `--quiet` suppresses informational output.
- `--verbose` enables verbose output (reserved for future diagnostics).

## Shell Integration

Use `--path-only` to auto-`cd` into the returned worktree path. Example for bash/zsh:

```bash
gtk() {
  local gittk
  gittk="$(command -v git-tk)" || { echo "git-tk not found in PATH"; return 127; }
  if [ "$1" = "branch" ] || [ "$1" = "checkout" ] || [ "$1" = "clone" ]; then
    local wt_path
    wt_path="$("$gittk" "$@" --path-only)"
    [ -d "$wt_path" ] && cd "$wt_path"
  else
    "$gittk" "$@"
  fi
}
```

## Examples

For detailed, scenario-based examples, see [docs/examples.md](docs/examples.md).

Quick reference:
```bash
git tk clone https://github.com/user/repo.git
cd repo/worktrees/main
git tk branch feature-x
git tk checkout feature-x
git tk list
git tk doctor
git tk prune
git tk sync --default
git tk setup --upstream-url https://github.com/org/repo.git
```

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
the first time you run it.

## License

MIT License. See [LICENSE](LICENSE) file.
