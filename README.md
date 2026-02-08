# ContextKeeper 🦞

Keep track of your thoughts, bugs, and ideas across all your devices without the hassle. Just a simple CLI that stores your project context where you already keep your code.

![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Platforms](https://img.shields.io/badge/platforms-Linux%20macOS%20Windows-lightgrey)

## What it does

- Stores context notes right in your project (`.contextkeeper/`) so they sync naturally with git
- Works on Linux, macOS, and Windows
- Minimal and fast - no startup lag, no bloat
- Tag and organize notes however you like

## Install

Grab binaries from the [releases page](https://github.com/ondrahracek/contextkeeper/releases):

**Linux:**
```bash
curl -L https://github.com/ondrahracek/contextkeeper/releases/download/v0.3.1/contextkeeper-linux-amd64.tar.gz -o ck.tar.gz
tar -xzf ck.tar.gz
chmod +x contextkeeper-linux-amd64
sudo mv contextkeeper-linux-amd64 /usr/local/bin/ck
```

**macOS:**
```bash
curl -L https://github.com/ondrahracek/contextkeeper/releases/download/v0.3.1/contextkeeper-darwin-arm64.tar.gz -o ck.tar.gz
tar -xzf ck.tar.gz
chmod +x contextkeeper-darwin-arm64
sudo mv contextkeeper-darwin-arm64 /usr/local/bin/ck
```

**Windows:**
```powershell
# Download and extract manually from releases page, or use:
Invoke-WebRequest -Uri "https://github.com/ondrahracek/contextkeeper/releases/download/v0.3.1/contextkeeper-windows-amd64.tar.gz" -OutFile ck.tar.gz
tar -xf ck.tar.gz
```

**Via Homebrew:**
```bash
brew install ondhrahracek/tap/contextkeeper
```

**From source:**
```bash
git clone https://github.com/ondrahracek/contextkeeper.git
cd contextkeeper
go build -o ck .
```

## Quick usage

Initialize a project:
```bash
ck init                    # Creates .contextkeeper/ in current directory
ck init --path ./shared   # Custom path for shared team context
ck init --global          # Global storage at ~/.local/share/contextkeeper
```

Add notes:
```bash
ck add "Remember to fix the auth bug in login.js"
ck add "API endpoint needs rate limiting" --project "api"
ck add "Great idea for feature X" --tags "idea,feature"
ck add -e                  # Opens your editor for longer notes
```

See what you have:
```bash
ck list                    # Show all active notes
ck list --project webapp   # Filter by project
ck list --tags bug         # Filter by tag
ck list --all              # Include completed items
ck list --json             # JSON output for scripting
```

Mark things done:
```bash
ck done <id>               # Mark item as completed
ck remove <id>             # Archive item
ck edit <id>               # Edit item content
```

## Where it stores things

ContextKeeper stores all your notes in a single file called `items.json` inside the `.contextkeeper/` directory. This file lives in your project and syncs naturally with git.

ContextKeeper looks for storage in this order:

1. Explicit path: `--path` flag
2. Environment variable: `CK_PATH`
3. Local project: `.contextkeeper/` directory
4. Global default: OS-specific location (e.g., `~/.local/share/contextkeeper`)

## Git sync

Since context lives in `.contextkeeper/`, it syncs naturally with git. Just add it to your repo:

```bash
git add .contextkeeper/
git commit -m "Add project context"
```

Or add to `.gitignore` if you prefer local-only storage.

## Commands at a glance

| Command | What it does |
|---------|--------------|
| `ck add [content]` | Add a new note |
| `ck list` | List all notes |
| `ck done <id>` | Mark as completed |
| `ck remove <id>` | Archive or delete |
| `ck edit <id>` | Edit a note |
| `ck init` | Set up storage |
| `ck status` | Quick overview |

## Building from source

```bash
git clone https://github.com/ondrahracek/contextkeeper.git
cd contextkeeper
go build -o ck .
```

Or use the Makefile:
```bash
make build              # Build for current platform
make build-all          # Build all platforms
make test               # Run tests
```

## Contributing

Found a bug? Have an idea? PRs are welcome!

1. Fork it
2. Create a feature branch
3. Make your changes
4. Open a PR

## Thanks to

- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Go](https://go.dev/) for making building tools this straightforward

---

**Never lose context again.** 🦞
