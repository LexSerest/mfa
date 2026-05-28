[![Go Report Card](https://goreportcard.com/badge/LexSerest/mfa)](https://goreportcard.com/report/LexSerest/mfa)
![GitHub Release](https://img.shields.io/github/v/release/LexSerest/mfa)
![GitHub License](https://img.shields.io/github/license/LexSerest/mfa)

# mfa

mfa is a fast, secure, and offline Command Line Interface (CLI) Two-Factor Authentication (2FA) manager written in Go. It keeps your tokens safe using industry-standard cryptography while providing instant access right from your terminal.

## Features

- **High Security:** Passwords derived via `PBKDF2-HMAC-SHA256` and database encrypted using `AES-GCM-256`.
- **Zero Dependencies:** Fully self-contained, working completely offline.
- **Smart Dual-Code Display:** Generates both the **current** and **next** TOTP codes simultaneously. Never rush again if a token is about to expire.
- **Terminal QR Codes:** Renders QR codes directly in your terminal for easy scanning with mobile authenticator apps.
- **Out-of-the-Box Autocompletion:** Tab-completion for `bash`, `zsh`, and `fish` sets up automatically upon installation.
- **Clipboard Integration:** Instantly copies the current active code to your clipboard.
---

## Installation & Autocompletion

When installed via package managers (`brew`, `nix`, `dpkg`, `rpm`), native shell tab-completion for `bash`, `zsh`, and `fish` is configured automatically and works out of the box.

### Using Nix (Flake)

You can run the application directly without installing it into your system:
```bash
nix run github:LexSerest/mfa -- --help
```

To install it into your NixOS configuration, add the repository to your inputs:
```nix
inputs.mfa-cli.url = "github:LexSerest/mfa";
```
Then include the package into your `environment.systemPackages`:
```nix
inputs.mfa-cli.packages.${pkgs.stdenv.hostPlatform.system}.default
```

### Using Homebrew (macOS / Linux)

```bash
brew tap LexSerest/tap
brew install mfa
```

### Linux Packages (.deb / .rpm)

Download the latest release for your architecture from the Releases page and install it using your package manager:

**Ubuntu / Debian:**
```bash
sudo dpkg -i mfa_amd64.deb
```

**Fedora / RHEL / CentOS:**
```bash
sudo rpm -i mfa_amd64.rpm
```

### From Source
Ensure you have Go 1.25+ installed. Note that manual source installation does not automatically copy completion files into system paths:
```bash
go install https://github.com/LexSerest/mfa
```

---

## Usage

```text
Usage:
  mfa <command> [arguments]

Commands:
  add         <label>               Add a new 2FA account [--digits=6] [--period=30] [--algo=SHA1]
  gen         <label>               Generate current TOTP code and copy to clipboard
  qr          <label>               Display account's QR code in the terminal for mobile sync [--big]
  list, ls                          List all saved account labels
  del, rm     <label>               Delete a 2FA account
  rename, mv  <label> <new label>   Rename an existing account label
  import      <path>                Import from plain text file (.txt with otpauth://)
  export      <path>                Export to UNENCRYPTED plain text file (otpauth:// format)
  help                              Show this help message
```

### Examples

**Adding a new token:**
```bash
mfa add google
# It will prompt for your secret key securely
```

**Generating a 2FA code:**
```bash
$ mfa gen google

875311 (22s) (copied!)
345666 (next)
```
*The current code is automatically pushed to your system clipboard for instant paste.*

**Syncing back with a phone app:**
```bash
# Displays a scannable QR code directly inside your terminal session 
# for easy secure synchronization with Aegis, Google Authenticator, etc.
mfa qr work-slack --big
```

**Renaming labels:**
```bash
mfa rename google google-personal
```

---

## Automation & Pipelines

The application fully supports non-interactive execution, scripting, and integration into terminal pipelines. It automatically detects whether `stdin` is a real terminal or a data pipe.

### Non-Interactive Password Input
You can pipe your master password directly into any command that requires authentication. This allows integration with system password managers like `pass` or `1Password CLI`:

```bash
# Using standard echo 
echo "my_master_password" | mfa gen google

# Using UNIX pass manager safely
pass show tools/mfa-master | mfa gen google
```

### Raw Token Output for Scripts
To use generated tokens inside automation scripts, use the `--raw` flag with the `gen` command. This disables the dual-code display, suppresses terminal formatting, skips clipboard copying, and prints **only** the pure 6 or 8-digit code string without trailing newlines:

```bash
# Store the active token inside a shell variable
TOKEN=$(echo "my_master_password" | mfa gen google --raw)

# Pass the raw token directly to another CLI utility or API request
pass show tools/mfa-master | mfa gen github --raw | xargs -I {} curl -H "Authorization: Bearer {}" https://api.github.com
```

---

## Backup, Import & Export

The manager allows seamless data migration between desktop and mobile devices using the standard OTP Auth URI format.

### Supported Format
* **Plain Text (otpauth://):** The application exclusively supports standard plain text files containing OTP Auth URIs. Applications like **Aegis Authenticator** can export backups directly into this format.

The backup file must contain exactly one valid URI per line. For example:
```text
otpauth://totp/Google:user?secret=KVKVE43V&issuer=Google
otpauth://totp/GitHub:developer?secret=KVKVE43V&issuer=GitHub
```

### Usage
To import or export your tokens, use the following commands:
```bash
mfa import backup.txt
mfa export backup.txt
```

---

## Database & Security

All your accounts are stored locally in an encrypted file. No telemetry, no cloud sync, complete privacy.

- **Storage Path:** `~/.config/mfa/` (or your OS equivalent platform storage path like AppData). Can be overridden via the `MFA_DB_PATH` environment variable.
- **Encryption:** Encryption keys are never stored. They are derived at runtime from your master password using **PBKDF2** with 10,000 iterations and a unique salt. The database itself is wrapped in **AES-GCM-256** authenticated encryption.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
