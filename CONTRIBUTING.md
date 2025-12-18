# Contributing

## Commit Message Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/) for automatic semantic versioning.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature (triggers minor version bump: 0.1.0 → 0.2.0)
- **fix**: A bug fix (triggers patch version bump: 0.1.0 → 0.1.1)
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **ci**: CI/CD configuration changes

### Breaking Changes

Add `BREAKING CHANGE:` in the commit body or footer to trigger a major version bump (0.1.0 → 1.0.0).

### Examples

**Feature (minor version bump):**

```
feat: add support for nested components

This allows components to contain other components as children.
```

**Bug fix (patch version bump):**

```
fix: resolve parsing error for self-closing tags
```

**Breaking change (major version bump):**

```
feat: redesign component API

BREAKING CHANGE: Component.Render() now requires a Registry parameter
```

**No version bump:**

```
docs: update README with installation instructions
```

```
chore: update dependencies
```

## Release Process

Releases are automated via GitHub Actions:

1. **Tests run on all branches** - Every push triggers the test workflow
2. **Releases created on main branch** - When pushing to `main`:
   - Tests are run
   - Version is calculated based on conventional commits since last tag
   - If version changed, a new release is created with binaries for:
     - Linux (AMD64, ARM64)
     - macOS (AMD64, ARM64/Apple Silicon)
     - Windows (AMD64)
   - Changelog is automatically generated
   - Binaries and SHA256 checksums are attached to the release

## Initial Version

The project starts at version **0.1.0**. The first push to main will create this initial release if no tags exist.

## Development Workflow

1. Create a feature branch: `git checkout -b feat/my-feature`
2. Make changes and commit using conventional commits: `git commit -m "feat: add new feature"`
3. Push to GitHub: `git push origin feat/my-feature`
4. Tests will run automatically
5. Create a pull request to `main`
6. After merge to `main`, a release will be automatically created if the version changed
