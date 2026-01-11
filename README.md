# gxcommit

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE.md)

A CLI that analyzes your git diff and generates logical git commits with AI using Groq. Automatically handles git add and git commit for you.

## Installation

### Install via Go

```sh
go install github.com/tahzeer/gxcommit@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

### Build from Source

Clone the repository and build:

```sh
git clone https://github.com/tahzeer/gxcommit.git
cd gxcommit
go build -o gxcommit
```

Move the binary to your PATH:

```sh
mv gxcommit /usr/local/bin/
# or
sudo mv gxcommit /usr/local/bin/
```

## Setup

1. Get your API key from [Groq Console](https://console.groq.com/keys)

2. Set the API key:

```sh
gxcommit config set GROQ_API_KEY=<your-api-key>
```

## Commands

### run

Generate and run commits immediately. Automatically handles git add and git commit based on the git diff.

```sh
gxcommit run [flags]
```

**Alias:** `r`

**Flags:**
- `-h, --help`: help for run

**Example:**
```sh
gxcommit run
```

### save-script

Generate a bash script with git add and git commit commands without executing them.

```sh
gxcommit save-script [flags]
```

**Alias:** `ss`

**Flags:**
- `-h, --help`: help for save-script

### config

Manage gxcommit configuration.

```sh
gxcommit config [command]
```

**Available Commands:**
- `set`: Set a configuration value

**Example:**
```sh
# Set API key
gxcommit config set GROQ_API_KEY=your-api-key
```

### version

Print the version number.

```sh
gxcommit version
```

### completion

Generate the autocompletion script for the specified shell.

```sh
gxcommit completion [bash|zsh|fish|powershell]
```

## Global Flags

These flags can be used with any command:

- `-c, --code string`: JIRA/ticket code to prefix commit messages

## Usage

Generate commits based on your git diff:

```sh
gxcommit run
```

With JIRA/ticket code prefix:

```sh
gxcommit run --code PROJ-123
# or
gxcommit -c PROJ-123 run
```

## Configuration

Manage configuration via the `config` command:

```sh
# Set API key
gxcommit config set GROQ_API_KEY=<your-api-key>
```

Configuration is stored in `~/.gxconfig` in gitconfig format:

```ini
[groq]
	GROQ_API_KEY=<your-api-key>
```

## Development

```sh
# Run locally
go run main.go

# Build
go build -o gxcommit

# Run tests
go test ./...
```

## License

MIT - see [LICENSE.md](LICENSE.md) file for details.
