# lnr — Linear CLI for humans & LLMs

A fast, focused CLI for [Linear](https://linear.app) that works equally well for interactive use and automation. Manage issues, cycles, projects, and initiatives from your terminal.

## Features

- **Issues** — list, view, create, update, close, archive, search, comment, label, and attach files
- **Cycles** — list, view current sprint, add/remove issues
- **Projects** — list, view details, list project issues
- **Initiatives** — list, view details, list initiative projects
- **Teams** — list teams, set a default team
- **Workflow states** — list available states for a team
- **Interactive dashboard** — TUI mode (`lnr tui`) with keyboard navigation
- **Plain output** — `--plain` flag for scripting and LLM-friendly output
- **Shell completions** — bash and zsh

## Installation

### From GitHub Releases

Download the latest binary from [Releases](https://github.com/misham/lnr/releases) and place it in your `$PATH`.

### From Source

Requires Go 1.22+.

```bash
git clone https://github.com/misham/lnr.git
cd lnr
make build
# binary is at bin/lnr
```

## Getting Started

```bash
# Authenticate with Linear (opens browser for OAuth)
lnr auth login

# Verify authentication
lnr auth status

# Set a default team so you don't need --team on every command
lnr team set
```

## Usage

### Issues

```bash
lnr issue list                    # List issues for your default team
lnr issue list --team ENG         # List issues for a specific team
lnr issue view ENG-123            # View issue details
lnr issue create --title "Bug"    # Create a new issue
lnr issue update ENG-123 --state "In Progress"
lnr issue close ENG-123
lnr issue search "login bug"
```

### Comments & Labels

```bash
lnr issue comment list ENG-123
lnr issue comment add ENG-123 --body "Looks good"
lnr issue label list ENG-123
lnr issue label add ENG-123 --label "bug"
```

### Cycles

```bash
lnr cycle list                    # List cycles
lnr cycle current                 # Show the active sprint
lnr cycle view <cycle-id>         # View cycle details
lnr cycle add-issue <issue-id>    # Add issue to current cycle
lnr cycle remove-issue <issue-id>
```

### Projects & Initiatives

```bash
lnr project list
lnr project view <project-id>
lnr project issues <project-id>

lnr initiative list
lnr initiative view <initiative-id>
lnr initiative projects <initiative-id>
```

### Interactive Dashboard

```bash
lnr tui
```

### Scripting & Automation

Use `--plain` to disable styled output for piping and LLM consumption:

```bash
lnr issue list --plain
lnr issue view ENG-123 --plain
```

### Shell Completions

```bash
# bash
lnr completion bash > /etc/bash_completion.d/lnr

# zsh
lnr completion zsh > "${fpath[1]}/_lnr"
```

## Development

```bash
make setup          # Install dev tools
make check          # Lint, vet, vulncheck
make test           # Run tests with race detector
make generate       # Regenerate GraphQL client from schema
make schema         # Fetch latest Linear schema
```
