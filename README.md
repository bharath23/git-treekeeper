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
- Listing worktrees
- Checking worktree health
- Pruning stale worktrees
- Syncing branches with remotes

More features like syncing, garbage collection, and stacked branches will be
added in future releases.

## CLI

The binary is named `git-tk`. Git discovers it as a plugin, so users invoke it
as `git tk ...`. Direct invocation with `git-tk ...` also works. Example usage:

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

Reports worktree health. Each branch is marked as `clean`, `dirty`, or
`merge/rebase in progress`.

```bash
git tk doctor
```

Options:
- `--porcelain` prints tab-separated rows with no header.
- `--json` prints a JSON array of `{branch,state}` objects.

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

### sync

Syncs a branch using `git fetch` plus `--ff-only` fast-forward, matching a
read-only main workflow.

```bash
git tk sync
git tk sync --default
git tk sync --branch feature-x
```

Options:
- `--default` syncs the default branch worktree.
- `--branch <name>` syncs the specified branch worktree.
- `--remote <name>` fetches from a specific remote.
- `--add-upstream <url>` adds an upstream remote (default name: `upstream`).
- `--set-upstream` sets the branch upstream to `<upstream>/<branch>` and sets
  push remote to `origin`.
- `--upstream <name>` sets the upstream remote name (default `upstream`).
- `--origin <name>` sets the origin remote name for push defaults.
- `--dry-run` prints what would be synced without making changes.

### Global Flags

- `--quiet` suppresses informational output.
- `--verbose` enables verbose output (reserved for future diagnostics).

## Shell Integration

Use `--path-only` to auto-`cd` into the returned worktree path. Example for bash/zsh:

```bash
gtk() {
  if [ "$1" = "branch" ] || [ "$1" = "checkout" ] || [ "$1" = "clone" ]; then
    local path
    path="$(command git tk "$@" --path-only)"
    [ -n "$path" ] && cd "$path"
  else
    command git tk "$@"
  fi
}
```

## Example Usage

```bash
git tk clone https://github.com/user/repo.git
cd repo/worktrees/main
git tk branch feature-x
git tk checkout feature-x
git tk list
git tk doctor
git tk prune
git tk sync --default
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
