# Git TreeKeeper: Example Driven Guide

This guide provides practical examples for common workflows using `git-treekeeper` (`git tk`).

## 1. Getting Started

### Cloning a Repository
Instead of a standard `git clone`, use `git tk clone`. This sets up a bare repository with a dedicated `worktrees/` directory.

```bash
# Clone the repo
git tk clone https://github.com/user/project.git

# Move into the default branch worktree (usually 'main' or 'master')
cd project/worktrees/main
```

### Setting up a Fork
If you've forked a repository, you can quickly configure it to track the original "upstream" repository.

```bash
# Add upstream and configure the current branch to track it
git tk setup --upstream-url https://github.com/original-owner/project.git --install-hooks
```
*The `--install-hooks` flag adds a pre-commit hook to prevent accidental commits to your main branch.*

---

## 2. Developing Features

### Creating a New Feature Branch
Create a new branch and a corresponding worktree in one command.

```bash
# Creates branch 'feature-abc' and worktree at 'worktrees/feature-abc'
git tk branch feature-abc

# If you use the shell integration (see README), you'll auto-cd into it.
# Otherwise:
cd ../feature-abc
```

### Switching Between Features
`git tk checkout` ensures a worktree exists for the branch and provides its path.

```bash
# Switch to an existing branch
cd $(git tk checkout existing-feature --path-only)

# If the branch doesn't exist, it creates it from the default branch
git tk checkout new-experiment
```

### Keeping Your Feature Up-to-Date
Sync your feature branch with its remote (origin).

```bash
# Performs a fetch and a fast-forward merge
git tk sync
```

---

## 3. Collaboration & Maintenance

### Syncing the Default Branch
Easily update your local `main` branch from the upstream repository.

```bash
# Sync the default branch worktree, even if you are in a feature worktree
git tk sync --default
```

### Checking Worktree Health
See which worktrees are clean, dirty, or have a rebase in progress.

```bash
git tk doctor
```
*Output:*
```text
branch      state                     tracking
------      -----                     --------
main        clean                     origin/main
feature-x   dirty                     origin/feature-x
feature-y   merge/rebase in progress  none
```

### Repairing Origin Tracking
If `git status` doesn't show ahead/behind in treekeeper-managed bare repos,
repair origin tracking and upstreams.

```bash
git tk repair
```

```bash
# Repair just one branch
git tk repair --branch feature-x --apply
```

```bash
# Explicitly fix tracking
git tk repair --tracking --apply --remote origin
```

```bash
# Apply changes (repair defaults to dry-run)
git tk repair --apply
```

### Listing All Worktrees
```bash
git tk list
```
*Output:*
```text
branch      path
------      ----
main        /path/to/project/worktrees/main
feature-x   /path/to/project/worktrees/feature-x
```

---

## 4. Cleanup

### Deleting a Merged Branch
Delete a branch and its worktree safely. `git tk` prevents deletion if the worktree is dirty or if the branch isn't merged into the default branch.

```bash
# Delete branch and worktree
git tk branch -d feature-done
```

### Force Deleting a Branch
If you want to abandon a feature that hasn't been merged:

```bash
# Use -D to force delete. --yes skips the confirmation prompt.
git tk branch -D feature-abandoned --yes
```

### Pruning Stale Worktrees
If you've manually deleted worktree directories, use `prune` to clean up Git's internal tracking.

```bash
# Preview what will be pruned
git tk prune --dry-run

# Perform the prune
git tk prune
```

### Garbage Collecting Old Branches
Clean up local branches that have been merged and are older than the threshold.
`gc` runs in dry-run mode by default and skips active worktrees.

```bash
# Preview what would be deleted (default)
git tk gc

# Delete branches older than 30 days (default threshold)
git tk gc --apply

# Use a custom threshold
git tk gc --apply --age-days 60
```

---

## 5. Advanced & Scripting

### Using JSON Output
For custom scripts or integration with other tools like `jq`.

```bash
# Get all dirty worktrees
git tk doctor --json | jq '.[] | select(.state == "dirty")'
```

### Seamless Navigation (Shell Integration)
Add this to your `.bashrc` or `.zshrc` to automatically `cd` into worktrees when using `branch`, `checkout`, or `clone`.

```bash
gtk() {
  local gittk
  gittk="$(command -v git-tk)" || { echo "git-tk not found in PATH"; return 127; }
  if [[ "$1" =~ ^(branch|checkout|clone)$ ]]; then
    local wt_path
    wt_path="$("$gittk" "$@" --path-only)"
    [ -d "$wt_path" ] && cd "$wt_path"
  else
    "$gittk" "$@"
  fi
}
---

## 6. Safety & Error Prevention

`git-tk` is designed to prevent accidental data loss. You might encounter errors in the following scenarios:

- **Dirty Worktree:** If you try to delete a branch whose worktree has uncommitted changes, `git-tk` will refuse. Commit or stash your changes first.
- **Unmerged Branch:** `git tk branch -d` will refuse to delete a branch that hasn't been merged into the default branch. Use `-D` to override this if you're sure.
- **Active Worktree:** You cannot delete a branch if your current shell is inside that branch's worktree. Move to a different directory first.
- **Merge/Rebase in Progress:** `git-tk doctor` will identify worktrees where a merge or rebase is stuck, and deletion will be blocked until the operation is completed or aborted.

---

## 7. Git Command Pass-Through


`git-tk` is designed to be a transparent addition to your Git workflow. If you run a command or flag that `git-tk` doesn't recognize, it automatically passes it through to the standard `git` command.

This means you can use `git tk` as a complete replacement for `git` in your terminal:

```bash
# Runs standard 'git status'
git tk status

# Runs standard 'git log --oneline'
git tk log --oneline

## 8. Shell Autocompletion

`git-tk` supports native tab completion for commands, branches, and remotes.

### Enabling Completion
To load completions for your current shell session:

**Bash:**
```bash
source <(git-tk completion bash)
```

**Zsh:**
```bash
source <(git-tk completion zsh)
```

**Fish:**
```bash
git-tk completion fish | source
```

### Wrapper Integration
If you use the `gtk` wrapper function, you must tell your shell to use `git-tk`'s completion logic for `gtk`.

**Zsh:**
Add this to your `.zshrc` *after* defining the function:
```bash
compdef gtk=git-tk
```

**Bash:**
Add this to your `.bashrc`:
```bash
complete -o default -F __start_git_tk gtk
```
