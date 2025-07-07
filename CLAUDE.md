# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

- Build the binary: `go build -o hatena-blog-org`
- Run tests: `go test -v`
- Run specific test: `go test -v -run TestFunctionName`

## Architecture

This is a Go CLI tool that converts Org-mode files to Markdown and posts them to Hatena Blog via AtomPub API.

### Core Components

- **main.go**: Entry point with CLI argument parsing and orchestration
- **config.go**: Configuration management (JSON config files, defaults to `~/.config/hatena-blog-org/config.json`)
- **converter.go**: Org-to-Markdown conversion using pandoc
- **hatena_client.go**: Hatena Blog API client with WSSE authentication

### Key Dependencies

- **pandoc**: External dependency for Org-to-Markdown conversion
- **Hatena Blog AtomPub API**: For posting entries

### Data Flow

1. Parse CLI arguments and load configuration
2. Convert .org file to Markdown using pandoc
3. Extract title from Markdown (first `# ` line)
4. Create XML entry with WSSE authentication
5. POST to Hatena Blog API

### Authentication

Uses WSSE (Web Services Security Username Token) authentication with SHA1 digest for Hatena Blog API access.