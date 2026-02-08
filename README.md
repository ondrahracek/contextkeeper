# ContextKeeper 🦞

A minimalist CLI tool for managing project context across multiple devices. Never lose track of your thoughts, bugs, or ideas while working on complex projects.

![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Platforms](https://img.shields.io/badge/Platforms-Linux%20%7C%20macOS%20%7C-Windows-lightgrey.svg)

## ✨ Features

- **Project-aware**: Automatically organizes context by project
- **Git-synced**: Store context in `.contextkeeper/` directory and sync with git
- **Cross-platform**: Linux, macOS, and Windows support
- **Minimalist**: No bloat, just what you need
- **Fast**: Written in Go, instant startup
- **Tag support**: Organize with tags like `bug`, `urgent`, `idea`

## 🚀 Quick Start

### Download Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/ondrahracek/contextkeeper/releases):

```bash
# Linux amd64
wget https://github.com/ondrahracek/contextkeeper/releases/download/v0.1.0/contextkeeper-linux-amd64
chmod +x contextkeeper-linux-amd64
sudo mv contextkeeper-linux-amd64 /usr/local/bin/ck

# macOS (Apple Silicon)
curl -L https://github.com/ondrahracek/contextkeeper/releases/download/v0.1.0/contextkeeper-darwin-arm64 -o ck
chmod +x ck
sudo mv ck /usr/local/bin/ck
```

### Install via Homebrew

```bash
brew install ondhrahracek/tap/contextkeeper
```

### Build from Source

```bash
git clone https://github.com/ondrahracek/contextkeeper.git
cd contextkeeper
go build -o ck .
sudo mv ck /usr/local/bin/ck
```

## 📖 Usage

### Initialize a Project

```bash
ck init                    # Creates .contextkeeper/ in current directory
ck init --path ./shared   # Custom path for team context
ck init --global          # Creates global context in ~/.local/share/contextkeeper
```

### Add Context Notes

```bash
ck add "Remember to fix the auth bug in login.js"
ck add "API endpoint needs rate limiting" --project "api"
ck add "Great idea for feature X" --tags "idea,feature"
ck add -e                  # Opens editor for multi-line notes
```

### List Your Context

```bash
ck list                    # Show all active notes
ck list --project webapp   # Filter by project
ck list --tags bug         # Filter by tag
ck list --all              # Include completed items
ck list --json             # JSON output for scripting
```

### Manage Items

```bash
ck done <id>               # Mark item as completed
ck remove <id>             # Archive item
ck remove <id> --force     # Permanently delete
ck edit <id>               # Edit item content
ck status                  # Quick overview
```

### Configuration

```bash
ck config --show           # Show current config
ck config --set default_project=myapp
ck config --get editor
```

## 📁 Storage Locations

ContextKeeper searches for storage in this order:

1. **Explicit path**: `--path` flag
2. **Environment variable**: `CK_PATH`
3. **Local project**: `.contextkeeper/` directory (git-synced)
4. **Global default**: OS-specific location (e.g., `~/.local/share/contextkeeper`)

### Git Integration

Add to your `.gitignore`:

```gitignore
.contextkeeper/
```

Or commit the `.contextkeeper/` directory to sync context across devices:

```bash
git add .contextkeeper/
git commit -m "Add project context"
```

## 🏗️ Architecture

```
contextkeeper/
├── cmd/
│   └── root.go           # CLI entry point
├── internal/
│   ├── cli/              # Command implementations
│   ├── config/           # Configuration and path detection
│   ├── models/           # Data structures
│   ├── storage/          # JSON persistence layer
│   └── utils/            # Helper functions
├── .contextkeeper/       # Default storage location
│   ├── config.json       # Configuration
│   └── data.json         # Context items
├── Makefile
└── README.md
```

## 🛠️ Development

### Build

```bash
make build              # Build for current platform
make build-all          # Build all platforms
make clean              # Clean build artifacts
```

### Test

```bash
make test               # Run all tests
make test-coverage      # Run tests with coverage
```

### Release

```bash
make release            # Create release binaries
```

## 📝 Commands Reference

| Command | Description |
|---------|-------------|
| `ck add [content]` | Add a new context note |
| `ck list` | List all context notes |
| `ck done <id>` | Mark item as completed |
| `ck remove <id>` | Archive or delete an item |
| `ck edit <id>` | Edit an item |
| `ck init` | Initialize storage directory |
| `ck config` | Manage configuration |
| `ck status` | Show quick overview |

### Global Flags

| Flag | Description |
|------|-------------|
| `--path, -p` | Custom storage path |
| `--help, -h` | Show help |
| `--version, -v` | Show version |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Go](https://go.dev/) for the excellent programming language
- All contributors and users!

---

**Never lose context again.** 🦞
