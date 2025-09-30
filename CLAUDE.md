# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a GitHub CLI extension (`gh-targetprocess`) that integrates Targetprocess with GitHub workflows. It creates pull requests using Targetprocess ticket information and allows viewing tickets from the command line.

## Key Commands

### Build
```bash
go build -o gh-targetprocess
```

### Install Extension Locally
```bash
gh extension install .
```

### Run
```bash
# Using gh CLI
gh targetprocess [ticket-id or url]

# Direct binary
./gh-targetprocess [ticket-id or url]
```

### View Ticket
```bash
gh targetprocess view [ticket-id] [--web]
```

### Configure
```bash
gh targetprocess configure
```

## Architecture

### Project Structure

- **main.go**: Entry point that initializes config (Viper), loads credentials from keyring, creates Targetprocess client, and executes root command
- **cmd/**: Cobra command definitions
  - **root.go**: Main command that creates PRs from Targetprocess tickets
  - **view/**: Subcommand to view ticket details (terminal or browser)
  - **configure/**: Subcommand to set up credentials
- **pkg/targetprocess/**: Targetprocess API client
  - **client.go**: HTTP client with API authentication (token via query param)
  - **models.go**: Assignable type and PR template generation
- **internal/config/**: Configuration management using Viper + go-keyring
  - Stores base URL in config file (`~/.config/gh-targetprocess.json`)
  - Stores access token in system keyring
- **internal/utils/**: Helper functions
  - **git.go**: Get current branch
  - **extraction.go**: Extract ticket ID from branch name or URL
- **templates/**: PR body template

### Key Design Patterns

**Context-based dependency injection**: Config and Targetprocess client are passed through context.Context using internal/context.go helpers (InitContext, GetConfig, GetTargetProcess)

**Ticket ID extraction**: Automatically extracts ID from:
1. Current git branch (pattern: `prefix/123_description`)
2. Targetprocess URL (pattern: `https://*.tpondemand.com/entity/123`)
3. Direct ID argument

**PR generation**: Fetches Assignable from Targetprocess API, converts HTML description to markdown, generates PR title/body, then calls `gh pr create`

**Configuration flow**:
- On first run, prompts for base URL and access token using charmbracelet/huh forms
- Validates token by testing API call to `/v1/Users/loggeduser`
- Stores token in system keyring (not in config file)

### API Client Notes

The Targetprocess client uses:
- `/v1/Assignables/{id}` to fetch ticket details
- Authentication via `access_token` query parameter
- Custom HTTP transport to add `Accept: application/json` header
- Deprecated: Direct httpClient field (use resty client instead)

### Important Files

- **cmd/root.go:63-65**: Main Assignable fetch logic
- **pkg/targetprocess/models.go:23-76**: PR title/body generation with HTML-to-markdown conversion
- **internal/utils/extraction.go:26-44**: Ticket ID extraction logic
- **internal/config/config.go:65-112**: Initial configuration setup with validation

## Dependencies

- **github.com/spf13/cobra**: CLI framework
- **github.com/cli/go-gh/v2**: GitHub CLI integration
- **github.com/spf13/viper**: Configuration management
- **github.com/zalando/go-keyring**: Secure credential storage
- **github.com/charmbracelet/huh**: Interactive forms
- **github.com/charmbracelet/glamour**: Markdown rendering
- **github.com/JohannesKaufmann/html-to-markdown/v2**: HTML conversion
- **resty.dev/v3**: HTTP client

## Release Process

GitHub Actions workflow triggers on:
- Manual dispatch
- Version tags (v*)
- Uses `cli/gh-extension-precompile@v2` for cross-platform builds