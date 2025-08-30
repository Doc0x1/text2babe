# 🔐 Text2Babe

A CLI encryption/encoding tool with an interactive shell interface, inspired by msfconsole. Supports AES-GCM encryption, multiple output formats, and Discord integration.

## Features

- **Interactive Shell**: msfconsole-style command interface with auto-completion
- **Dual Mode Operation**: Toggle between encryption (AES-GCM) and plain encoding
- **Multiple Output Formats**: hex, base64, or binary representation
- **Discord Integration**: Send encrypted messages directly to Discord DMs
- **Smart Auto-Detection**: Automatically detects input format when decrypting
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Clipboard Integration**: Automatically copies results to clipboard

## Installation

```bash
# Clone the repository
git clone github.com/doc0x1/text2babe
cd text2babe

# Install dependencies
go mod tidy

# Build for current platform
go build -o text2babe main.go

# Cross-platform builds

# Windows (Command Prompt)
set GOOS=windows&& set GOARCH=amd64&& go build -o text2babe.exe main.go

# Windows (PowerShell)
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o text2babe.exe main.go

# macOS/Linux - Build for Windows
GOOS=windows GOARCH=amd64 go build -o text2babe.exe main.go

# macOS/Linux - Build for macOS Intel
GOOS=darwin GOARCH=amd64 go build -o text2babe-macos-intel main.go

# macOS/Linux - Build for macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o text2babe-macos-arm main.go

# macOS/Linux - Build for Linux
GOOS=linux GOARCH=amd64 go build -o text2babe-linux main.go
```

## Quick Start

```bash
# Run interactive shell (default)
./text2babe.exe

# Direct command usage
./text2babe.exe encrypt "hello world"
./text2babe.exe decrypt <encrypted_data>
```

## Interactive Shell Commands

| Command | Description |
|---------|-------------|
| `encrypt <data>` | Encrypt/encode text data |
| `decrypt <data>` | Decrypt/decode data (auto-detects format) |
| `mode [encrypt/decrypt]` | Set or show current mode |
| `key <password>` | Set encryption key from password |
| `set <setting> <value>` | Configure settings |
| `toggle <setting>` | Toggle settings on/off |
| `discord [test/fetch]` | Discord operations |
| `config` | Show current configuration |
| `help` | Show available commands |

## Settings

| Setting | Values | Description |
|---------|--------|-------------|
| `encryption` | on/off | Enable/disable AES-GCM encryption |
| `output` | hex/base64/binary | Output format for encrypted data |
| `discord` | on/off | Auto-send to Discord DM |
| `discord-id` | channel_id | Set Discord DM channel ID |

## Examples

```bash
# Basic encryption
encrypt hello world

# Change output format
set output base64
encrypt secret message

# Toggle encryption off for plain encoding
set encryption off
encrypt plaintext    # Just converts to hex/base64/binary

# Discord operations
set discord-id 123456789012345678
set discord on
encrypt confidential  # Sends to Discord automatically
discord fetch         # Fetch and decrypt last Discord message
```

## Discord Integration

### Setup

1. Copy `.env.example` to `.env`
2. Set your Discord user token and DM channel ID
3. **Warning**: Using user tokens violates Discord's ToS - use at your own risk

### Features

- **Auto-Send**: Automatically send encrypted data to Discord DMs
- **Message Fetching**: Retrieve and decrypt the last text2babe message
- **Format Support**: Works with hex, base64, and binary formats
- **Validation**: Ensures fetched messages are from text2babe

## Security Features

- **AES-256-GCM**: Modern authenticated encryption
- **Key Derivation**: SHA-256 based key generation
- **No History**: Commands are not saved to disk
- **Auto-Detection**: Smart format detection prevents data corruption

## Configuration

The tool supports both environment variables and in-app configuration:

```bash
# Set via environment (.env file)
DISCORD_USER_TOKEN=your_token_here
DISCORD_DM_ID=channel_id_here

# Set via commands
set discord-id 123456789012345678
key mypassword
```

## Architecture

- **cmd/**: Cobra command definitions and interactive shell
- **internal/config/**: Configuration management
- **internal/crypto/**: AES-GCM encryption implementation  
- **internal/style/**: Cross-platform terminal styling
- **internal/discord/**: Discord API integration
- **pkg/prompt/**: Readline-based terminal interface

## License

This project is for educational purposes. Use responsibly and in compliance with applicable laws and terms of service.
