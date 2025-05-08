# VCS - Version Control System

**VCS** is an opinionated Git clone written in Go. It is a minimal, educational version control system designed to help developers understand how Git works at the object level — from hashing and content storage to plumbing commands like `cat-file`.

> ⚠️ This project is for learning purposes only. It is not intended for production use.

---

## ✅ Current Features

### `vcs init`
Initializes a new VCS repository by creating the `.vcs/` directory structure:
- `.vcs/objects/` — stores content-addressable blob objects
- `.vcs/refs/` — placeholder for future branch references
- `.vcs/HEAD` — points to the default branch (e.g., `main`)

### `vcs add <file>`
Simulates `git hash-object -w`:
- Reads the file content
- Prepares a blob object: `blob <size>\0<content>`
- Hashes the blob using **SHA-1**
- Compresses it with **zlib**
- Stores it under `.vcs/objects/<first-two-hash-chars>/<rest>`

### `vcs cat-file -p <hash>`
Replicates `git cat-file -p`:
- Locates and decompresses the stored blob
- Prints out the raw file content (excluding metadata)

### `.vcsignore` support (like `.gitignore`)
Ignores files and directories listed in a `.vcsignore` file during the `add` process:
- Glob-style patterns supported (e.g., `*.log`, `node_modules/`)
- Prevents accidental tracking of unwanted or system-specific files

---

## 🛠️ Upcoming Features

### `vcs status`
- Will show tracked and untracked changes
- Detect file modifications using SHA-1 comparisons

### `vcs commit`
- Will snapshot the current state of the index (added files)
- Add commit metadata like author, timestamp, and message

### `vcs cherry-pick <commit>`
- Will apply changes from a specified commit onto the current HEAD

---

## 🧱 Project Goals

- Replicate key Git internals from scratch
- Learn how Git manages objects, trees, and refs
- Build an extensible CLI using Go and Cobra
- Enable experimentation with version control logic

---

## 🧪 Getting Started

```bash
# Build the binary
go build -o vcs

# Initialize a VCS repo
./vcs init

# Add a file
echo "hello world" > test.txt
./vcs add test.txt

# Ignore some files
echo "*.log" > .vcsignore

# Read object content from hash
./vcs cat-file -p <hash-of-test.txt>
