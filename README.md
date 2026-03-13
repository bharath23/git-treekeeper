![](https://github.com/bharath23/git-treekeeper/actions/workflows/test-suite.yml/badge.svg)

# Git TreeKeeper

**Git TreeKeeper** is a CLI tool to help manage Git branches and worktrees
efficiently.

## Overview

For now, Git TreeKeeper focuses on:

- Cloning repositories
- Creating new branches and worktrees
- Checking out branches

More features like branch deletion, syncing, garbage collection, and stacked
branches will be added in future releases.

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

```bash
git tk branch <branch-name> [base]
```

### checkout

Ensures a worktree exists for the branch and prints its path. If the branch
does not exist, it creates it from the default branch.

```bash
git tk checkout <branch-name>
```

### clone

Clones a repository into the default directory (same as `git clone`), creates a
bare repo at `<base>/repo.git`, and adds a worktree for the default branch at
`<base>/worktrees/<branch>`.

```bash
git tk clone <repo-url> [path]
```

## Example Usage

```bash
git tk clone https://github.com/user/repo.git
cd repo/worktrees/main
git tk branch feature-x
git tk checkout feature-x
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

## License

MIT License. See [LICENSE](LICENSE) file.
