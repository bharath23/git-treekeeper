# Command Reference

This document describes each command and its flags in detail.

## branch

Creates a new branch and a worktree for it under `worktrees/<branch>`. If no
base is provided, it uses the current branch and falls back to the default
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

## checkout

Ensures a worktree exists for the branch and prints its path. If the branch
does not exist, it creates it from the default branch.

```bash
git tk checkout <branch-name>
```

Options:
- `--path-only` prints only the worktree path.

## clone

Clones a repository into the default directory (same as `git clone`), creates a
bare repo at `<base>/repo.git`, and adds a worktree for the default branch at
`<base>/worktrees/<branch>`.

```bash
git tk clone <repo-url> [path]
```

Options:
- `--path-only` prints only the worktree path.

## repair

Repairs tracking for treekeeper-managed bare repos. This configures the remote
fetch refspec, fetches the remote, and sets missing branch upstreams when the
remote branch exists.

```bash
git tk repair
```

Options:
- `--branch <name>` only repairs the specified branch.
- `--dry-run` shows what would change without making updates (default).
- `--apply` applies changes (disables dry-run).
- `--tracking` fixes tracking (default).
- `--remote <name>` sets which remote to track (default `origin`).

## doctor

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

## list

Lists all worktrees with their branch name and path, sorted by branch name.

```bash
git tk list
```

Options:
- `--porcelain` prints tab-separated rows with no header.
- `--json` prints a JSON array of `{branch,path}` objects.

## prune

Prunes worktree entries whose directories no longer exist. By default it only
removes stale worktrees; use `--merged-branches` to also remove merged branches
that no longer have a worktree.

```bash
git tk prune
```

Options:
- `--dry-run` prints what would be pruned without making changes.
- `--merged-branches` also deletes merged branches that have no worktree.

## gc

Garbage collects local branches that are merged and older than the threshold.
Runs in dry-run mode by default and skips active worktrees.

```bash
git tk gc
git tk gc --apply
```

Options:
- `--apply` deletes branches (default is dry-run).
- `--age-days <days>` only considers branches older than the threshold (default 30).

## sync

Syncs a branch using `git fetch` plus `--ff-only` fast-forward, matching a
read-only main workflow. By default, it syncs only clean worktrees.

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

## setup

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

## Global Flags

- `--quiet` suppresses informational output.
- `--verbose` enables verbose output (reserved for future diagnostics).
