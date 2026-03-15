#!/usr/bin/env bash
set -euo pipefail

BIN="${GIT_TK_BIN:-./build/git-tk}"
if [ ! -x "$BIN" ]; then
  echo "git-tk binary not found or not executable: $BIN" >&2
  exit 1
fi

BIN="$(cd "$(dirname "$BIN")" && pwd)/$(basename "$BIN")"

tmp="$(mktemp -d)"
cleanup() {
  rm -rf "$tmp"
}
trap cleanup EXIT

assert_contains() {
  local haystack="$1"
  local needle="$2"
  if ! printf "%s" "$haystack" | grep -Fq "$needle"; then
    echo "Expected output to contain: $needle" >&2
    echo "Actual output:" >&2
    printf "%s\n" "$haystack" >&2
    exit 1
  fi
}

assert_contains_any() {
  local haystack="$1"
  shift
  for needle in "$@"; do
    if printf "%s" "$haystack" | grep -Fq "$needle"; then
      return 0
    fi
  done
  echo "Expected output to contain one of: $*" >&2
  echo "Actual output:" >&2
  printf "%s\n" "$haystack" >&2
  exit 1
}

assert_contains_tabbed_path() {
  local haystack="$1"
  local branch="$2"
  local path="$3"
  local normalized
  normalized="$(normalize_path "$path")"
  if printf "%s" "$haystack" | awk -v b="$branch" -v p="$path" -v n="$normalized" '
    $1 == b && ($2 == p || $2 == n) { found=1 }
    END { exit found ? 0 : 1 }
  '; then
    return 0
  fi
  echo "Expected output to contain branch/path: $branch <tab> $path" >&2
  echo "Actual output:" >&2
  printf "%s\n" "$haystack" >&2
  exit 1
}

normalize_path() {
  local path="$1"
  if [ -d "$path" ]; then
    (cd "$path" && pwd -P)
    return
  fi
  local dir
  dir="$(dirname "$path")"
  local base
  base="$(basename "$path")"
  (cd "$dir" && printf "%s/%s" "$(pwd -P)" "$base")
}

assert_contains_path() {
  local haystack="$1"
  local path="$2"
  local normalized
  normalized="$(normalize_path "$path")"
  assert_contains_any "$haystack" "$path" "$normalized"
}

expect_fail() {
  local desc="$1"
  shift
  echo
  echo "== $desc =="
  if output="$("$@" 2>&1)"; then
    echo "Expected failure but command succeeded" >&2
    echo "Output:" >&2
    printf "%s\n" "$output" >&2
    exit 1
  fi
  echo "$output"
}

echo "Temp dir: $tmp"

echo
echo "== Setup source repo =="
mkdir -p "$tmp/src"
git -C "$tmp/src" init -b main
git -C "$tmp/src" config user.name git-tk
git -C "$tmp/src" config user.email git-tk@example.com
echo "hello" > "$tmp/src/README.md"
git -C "$tmp/src" add README.md
git -C "$tmp/src" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "init"

cd "$tmp"

echo
echo "== Clone (default destination) =="
output="$("$BIN" clone "$tmp/src" 2>&1)"
echo "$output"
assert_contains "$output" "Cloning repo $tmp/src"
assert_contains "$output" "Default branch: main"
assert_contains_path "$output" "$tmp/src/worktrees/main"
[ -d "$tmp/src/repo.git" ] || { echo "repo.git missing" >&2; exit 1; }
[ -d "$tmp/src/worktrees/main" ] || { echo "main worktree missing" >&2; exit 1; }
git -C "$tmp/src/worktrees/main" config user.name git-tk
git -C "$tmp/src/worktrees/main" config user.email git-tk@example.com

echo
echo "== Clone (explicit destination) =="
output="$("$BIN" clone "$tmp/src" "$tmp/myrepo" 2>&1)"
echo "$output"
assert_contains_path "$output" "$tmp/myrepo/worktrees/main"
[ -d "$tmp/myrepo/repo.git" ] || { echo "repo.git missing" >&2; exit 1; }

cd "$tmp/src/worktrees/main"

echo
echo "== Branch creation =="
output="$("$BIN" branch feature-x 2>&1)"
echo "$output"
assert_contains "$output" "Creating branch feature-x from main"
assert_contains_path "$output" "$tmp/src/worktrees/feature-x"

echo
echo "== Checkout existing branch =="
output="$("$BIN" checkout feature-x 2>&1)"
echo "$output"
assert_contains_path "$output" "$tmp/src/worktrees/feature-x"

echo
echo "== Checkout new branch =="
output="$("$BIN" checkout feature-y 2>&1)"
echo "$output"
assert_contains_path "$output" "$tmp/src/worktrees/feature-y"

echo
echo "== Checkout (path-only) =="
path="$("$BIN" checkout --path-only feature-path 2>&1)"
echo "$path"
assert_contains_path "$path" "$tmp/src/worktrees/feature-path"
cd "$path"
cd "$tmp/src/worktrees/main"

echo
echo "== List worktrees =="
output="$("$BIN" list 2>&1)"
echo "$output"
assert_contains "$output" "branch"
assert_contains "$output" "path"
assert_contains_path "$output" "$tmp/src/worktrees/feature-x"

echo
echo "== List (porcelain) =="
output="$("$BIN" list --porcelain 2>&1)"
echo "$output"
assert_contains_tabbed_path "$output" "main" "$tmp/src/worktrees/main"

echo
echo "== List (json) =="
output="$("$BIN" list --json 2>&1)"
echo "$output"
assert_contains "$output" "\"branch\":\"main\""
assert_contains_any "$output" "\"path\":\"$tmp/src/worktrees/main\"" "\"path\":\"$(normalize_path "$tmp/src/worktrees/main")\""

echo
echo "== Doctor (dirty) =="
echo "dirty" >> "$tmp/src/worktrees/main/dirty.txt"
output="$("$BIN" doctor 2>&1)"
echo "$output"
assert_contains "$output" "branch"
assert_contains "$output" "state"
assert_contains "$output" "dirty"

echo
echo "== Doctor (porcelain) =="
output="$("$BIN" doctor --porcelain 2>&1)"
echo "$output"
assert_contains_tabbed_path "$output" "main" "dirty"

echo
echo "== Doctor (json) =="
output="$("$BIN" doctor --json 2>&1)"
echo "$output"
assert_contains "$output" "\"branch\":\"main\""
assert_contains "$output" "\"state\":\"dirty\""

echo
echo "== Prune stale worktree =="
output="$("$BIN" branch stale-branch 2>&1)"
echo "$output"
rm -rf "$tmp/src/worktrees/stale-branch"
output="$("$BIN" prune 2>&1)"
echo "$output"
assert_contains_path "$output" "$tmp/src/worktrees/stale-branch"

echo
echo "== Prune (dry run) =="
output="$("$BIN" branch dry-branch 2>&1)"
echo "$output"
rm -rf "$tmp/src/worktrees/dry-branch"
output="$("$BIN" prune --dry-run 2>&1)"
echo "$output"
assert_contains "$output" "Would prune worktree:"
assert_contains_path "$output" "$tmp/src/worktrees/dry-branch"

echo
echo "== Prune merged branches =="
git -C "$tmp/src/worktrees/main" checkout -b merged-branch
echo "merged" >> "$tmp/src/worktrees/main/merged.txt"
git -C "$tmp/src/worktrees/main" add merged.txt
git -C "$tmp/src/worktrees/main" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "merge me"
git -C "$tmp/src/worktrees/main" checkout main
git -C "$tmp/src/worktrees/main" merge merged-branch
output="$("$BIN" prune --merged-branches 2>&1)"
echo "$output"
assert_contains "$output" "Pruned branch: merged-branch"

echo
echo "== Sync default from upstream =="
git clone --bare "$tmp/src/repo.git" "$tmp/upstream.git" >/dev/null
output="$("$BIN" sync --default --add-upstream "$tmp/upstream.git" --set-upstream 2>&1)"
echo "$output"
git clone "$tmp/upstream.git" "$tmp/upstream-work" >/dev/null
git -C "$tmp/upstream-work" config user.name git-tk
git -C "$tmp/upstream-work" config user.email git-tk@example.com
echo "upstream change" >> "$tmp/upstream-work/README.md"
git -C "$tmp/upstream-work" add README.md
git -C "$tmp/upstream-work" commit -m "upstream change"
git -C "$tmp/upstream-work" push origin main >/dev/null
rm -rf "$tmp/upstream-work"
output="$("$BIN" sync --default 2>&1)"
echo "$output"
assert_contains "$output" "Syncing main from upstream/main"

echo
echo "== Branch delete (unmerged) =="
output="$("$BIN" branch feature-del 2>&1)"
echo "$output"
echo "change" >> "$tmp/src/worktrees/feature-del/feature.txt"
git -C "$tmp/src/worktrees/feature-del" add feature.txt
git -C "$tmp/src/worktrees/feature-del" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "feature"
expect_fail "Delete unmerged branch" "$BIN" branch -d feature-del

echo
echo "== Branch delete (merged) =="
output="$("$BIN" branch feature-merged 2>&1)"
echo "$output"
echo "merged" >> "$tmp/src/worktrees/feature-merged/merged.txt"
git -C "$tmp/src/worktrees/feature-merged" add merged.txt
git -C "$tmp/src/worktrees/feature-merged" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "merge me"
git -C "$tmp/src/worktrees/main" merge feature-merged
output="$("$BIN" branch -d feature-merged 2>&1)"
echo "$output"
assert_contains "$output" "Deleted branch: feature-merged"

echo
echo "== Branch delete (force) =="
output="$("$BIN" branch -D --yes feature-del 2>&1)"
echo "$output"
assert_contains "$output" "Deleted branch: feature-del"

echo
echo "== Branch delete (dirty) =="
output="$("$BIN" branch feature-dirty 2>&1)"
echo "$output"
echo "dirty" >> "$tmp/src/worktrees/feature-dirty/dirty.txt"
expect_fail "Delete dirty branch" "$BIN" branch -d feature-dirty

echo
echo "== Branch delete (merge in progress) =="
output="$("$BIN" branch feature-merge 2>&1)"
echo "$output"
echo "feature change" >> "$tmp/src/worktrees/feature-merge/feature.txt"
git -C "$tmp/src/worktrees/feature-merge" add feature.txt
git -C "$tmp/src/worktrees/feature-merge" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "feature change"
echo "main change" >> "$tmp/src/worktrees/main/main.txt"
git -C "$tmp/src/worktrees/main" add main.txt
git -C "$tmp/src/worktrees/main" -c user.name=git-tk -c user.email=git-tk@example.com commit -m "main change"
git -C "$tmp/src/worktrees/feature-merge" merge main --no-ff --no-commit
expect_fail "Delete branch with merge in progress" "$BIN" branch -d feature-merge
git -C "$tmp/src/worktrees/feature-merge" merge --abort

echo
echo "== Branch delete (remote) =="
git init --bare -b main "$tmp/remote.git" >/dev/null
git -C "$tmp/src" remote add origin "$tmp/remote.git"
git -C "$tmp/src" push -u origin main
git -C "$tmp/remote.git" symbolic-ref HEAD refs/heads/main
output="$("$BIN" clone "$tmp/remote.git" "$tmp/remote-clone" 2>&1)"
echo "$output"
cd "$tmp/remote-clone/worktrees/main"
output="$("$BIN" branch feature-remote 2>&1)"
echo "$output"
git -C "$tmp/remote-clone/worktrees/feature-remote" push -u origin feature-remote
output="$("$BIN" branch -d --remote --yes feature-remote 2>&1)"
echo "$output"
assert_contains "$output" "Deleted remote branch: origin/feature-remote"

echo
echo "== Too many args =="
expect_fail "Clone too many args" "$BIN" clone "$tmp/src" one two
expect_fail "Branch too many args" "$BIN" branch a b c
expect_fail "Checkout too many args" "$BIN" checkout a b

echo
echo "== Help formatting =="
output="$("$BIN" -h 2>&1)"
echo "$output"
assert_contains "$output" "Usage:"
output="$("$BIN" branch -h 2>&1)"
echo "$output"
assert_contains "$output" "Usage:"

echo
echo "== Quiet output =="
output="$("$BIN" --quiet branch quiet-test 2>&1)"
if [ -n "$output" ]; then
  echo "Expected no output for --quiet, got:" >&2
  printf "%s\n" "$output" >&2
  exit 1
fi

echo
echo "All manual tests passed."
