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

The CLI binary will be named `git-tk`. Example usage:

```bash
git tk clone <repo-url>
git tk branch <branch-name>
git tk checkout <branch-name>
```

## License

MIT License. See [LICENSE](LICENSE) file.
