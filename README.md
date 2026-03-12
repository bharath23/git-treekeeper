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

## CLI (Initial Stage)

The binary is named `git-tk`. Git discovers it as a plugin, so users invoke it
as `git tk ...`. Direct invocation with `git-tk ...` also works. Example usage:

```bash
git tk clone <repo-url>
git tk branch <branch-name>
git tk checkout <branch-name>
```

## License

MIT License. See [LICENSE](LICENSE) file.
