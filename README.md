# h1 — HackerOne CLI

A Go CLI tool wrapping the [HackerOne Customer Resources REST API](https://api.hackerone.com/). Provides programmatic access to all stable HackerOne API functionality via subcommands. Supports JSON, plaintext, and Markdown output formats.

## Installation

### From source

```bash
task build
```

The binary is placed at `bin/h1`.

### Go install

```bash
go install github.com/scottbrown/hackerone-cli/cmd/h1@latest
```

## Authentication

The CLI uses HTTP Basic Auth with your HackerOne API credentials. Set these environment variables:

```bash
export HACKERONE_API_IDENTIFIER="your-api-identifier"
export HACKERONE_API_TOKEN="your-api-token"
```

You can also pass them as flags:

```bash
h1 --api-identifier "your-id" --api-token "your-token" reports list
```

Credentials can be created in your HackerOne organization settings under **API Tokens**.

## Usage

```bash
h1 [command] [subcommand] [flags]
```

Output defaults to JSON for easy scripting and piping to `jq`. Use the `--format` flag to switch between output formats. Errors are written to stderr.

### Output formats

| Format | Description |
|--------|-------------|
| `json` | Pretty-printed JSON (default) |
| `text` | Human-readable plaintext — key-value pairs for single items, tab-aligned tables for lists |
| `markdown` | Markdown tables suitable for pasting into issues or documents |

```bash
# Default JSON output
h1 reports list --program koho

# Human-readable text table
h1 reports list --program koho --format text

# Markdown table
h1 reports list --program koho --format markdown

# Single item as key-value pairs
h1 reports get 12345 --format text
```

### Quick examples

```bash
# List your programs
h1 programs list

# Get a specific report
h1 reports get 12345

# List reports with pagination
h1 reports list --page 1 --page-size 25

# Add a comment to a report
h1 reports comment 12345 --message "Triaged and confirmed." --internal

# Change report state
h1 reports state 12345 --state triaged --message "Confirmed vulnerability."

# Award a bounty
h1 reports bounties award 12345 --amount 500 --message "Thanks for the report!"

# List organization members
h1 organizations members list org-123

# Get analytics
h1 analytics get --program my-program --start-date 2024-01-01 --end-date 2024-12-31

# Send an email
h1 email send --to user@example.com --subject "Test" --body "Hello"

# Pipe to jq
h1 programs list | jq '.[].attributes.handle'
```

## Commands

| Command | Description |
|---------|-------------|
| `h1 activities` | Manage activities (get, list) |
| `h1 analytics` | Query analytics data |
| `h1 assets` | Manage assets, ports, scopes, tags, scanner config |
| `h1 automations` | Manage automations and runs |
| `h1 credentials` | Manage credentials, inquiries, and responses |
| `h1 email` | Send email via HackerOne |
| `h1 organizations` | Manage organizations, groups, members, invitations |
| `h1 programs` | Manage programs, scopes, bounties, invitations, policy |
| `h1 reports` | Manage reports, comments, state, bounties, attachments |
| `h1 users` | Look up users by username or ID |
| `h1 version` | Print version information |

Use `h1 [command] --help` to see available subcommands and flags.

## Development

### Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev/) (Go Task runner)

### Build

```bash
task build
```

### Test

```bash
task test
```

### Lint

```bash
task lint
```

### Clean

```bash
task clean
```

## License

See [LICENSE](LICENSE) for details.
