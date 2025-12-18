# staticgen

Generates assets which can be easily hosted via normal web servers. Especially useful if you're using something like htmx.

## Project Structure

```
staticgen/
├── pkg/                    # Main package code
│   ├── component/         # Component definitions and instances
│   ├── page/             # Page loading and generation
│   ├── htmlutil/         # HTML utilities
│   └── internal/         # Internal utilities (xmlutil)
├── tests/                # End-to-end tests
│   └── fixtures/        # Test fixtures
├── examples/            # Example projects
├── main.go             # CLI entry point
└── Makefile            # Build automation
```

## Prerequisites

This project requires `make` to build.

### Installing Make

**Linux:**

```bash
# Debian/Ubuntu
sudo apt-get install make

# Fedora/RHEL
sudo dnf install make

# Arch Linux
sudo pacman -S make
```

**Windows:**

```powershell
# Using Chocolatey
choco install make

# Or using Scoop
scoop install make

# Or install via WSL (Windows Subsystem for Linux)
wsl --install
```

## Building

```bash
make build          # Build the binary
make test           # Run tests
make run            # Build and run
make all            # Clean, test, and build
make help           # Show all commands
```

## Contributing

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated semantic versioning and releases. See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Documentation

- [USING.md](USING.md) - How to use staticgen with examples
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development and contribution guidelines
