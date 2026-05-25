package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type CLIParser struct {
	Command  string
	Label    string
	NewLabel string
	Path     string
	Digits   int
	Period   int
	Algo     string
	IsBig    bool
	IsRaw    bool
}

func parseArgs() *CLIParser {
	if len(os.Args) < 2 {
		return &CLIParser{Command: "help"}
	}

	cmd := os.Args[1]

	switch cmd {
	case "ls":
		cmd = "list"
	case "rm":
		cmd = "del"
	case "mv":
		cmd = "rename"
	}

	if cmd == "-h" || cmd == "--help" || cmd == "help" {
		return &CLIParser{Command: "help"}
	}

	cli := &CLIParser{Command: cmd}
	fs := pflag.NewFlagSet(cmd, pflag.ContinueOnError)

	fs.SetInterspersed(true)
	fs.SetOutput(os.Stderr)

	switch cmd {
	case "add":
		fs.IntVar(&cli.Digits, "digits", 6, "number of digits in TOTP code")
		fs.IntVar(&cli.Period, "period", 30, "TOTP step period in seconds")
		fs.StringVar(&cli.Algo, "algo", "SHA1", "hashing algorithm (SHA1, SHA256, SHA512)")
	case "gen":
		fs.BoolVar(&cli.IsRaw, "raw", false, "output only the raw current token")
	case "qr":
		fs.BoolVar(&cli.IsBig, "big", false, "render larger QR code format")
	}

	err := fs.Parse(os.Args[2:])
	if err != nil {
		os.Exit(1)
	}
	args := fs.Args()

	switch cmd {
	case "rename":
		if len(args) >= 1 {
			cli.Label = strings.TrimSpace(args[0])
		}
		if len(args) >= 2 {
			cli.NewLabel = strings.TrimSpace(args[1])
		}
	case "import", "export":
		if len(args) >= 1 {
			cli.Path = strings.TrimSpace(args[0])
		}
	case "completion":
		if len(args) >= 1 {
			cli.Path = strings.TrimSpace(args[0])
		}
	default:
		if len(args) >= 1 {
			cli.Label = strings.TrimSpace(args[0])
		}
	}

	switch cmd {
	case "add", "gen", "qr", "del":
		if cli.Label == "" {
			fmt.Fprintln(os.Stderr, "Error: label is required")
			os.Exit(1)
		}
	case "rename":
		if cli.Label == "" || cli.NewLabel == "" {
			fmt.Fprintln(os.Stderr, "Usage: mfa rename <old_name> <new_name>")
			os.Exit(1)
		}
	case "import":
		if cli.Path == "" {
			fmt.Fprintln(os.Stderr, "Usage: mfa import <file_path>")
			os.Exit(1)
		}
	case "export":
		if cli.Path == "" {
			fmt.Fprintln(os.Stderr, "Usage: mfa export <file_path>")
			os.Exit(1)
		}
	case "completion":
		if cli.Path != "bash" && cli.Path != "zsh" && cli.Path != "fish" {
			fmt.Fprintln(os.Stderr, "Usage: mfa completion [bash | zsh | fish]")
			os.Exit(1)
		}
	}
	return cli
}
