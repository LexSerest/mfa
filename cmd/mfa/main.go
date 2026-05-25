package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"

	"mfa/assets"
	"mfa/internal/models"
	"mfa/internal/storage"
	"mfa/internal/term"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func help() {
	fmt.Printf(assets.HelpText, storage.GetDbPath())
}

var Version = "development"

func main() {
	cli := parseArgs()

	switch cli.Command {
	case "version":
		fmt.Printf("mfa version %s\n", Version)
		return

	case "add":
		password := term.TermGetPassword("Password: ")
		secret := term.TermGetPassword(fmt.Sprintf("2fa secret '%s': ", cli.Label))
		account := models.NewAccount(secret, cli.Digits, cli.Period, cli.Algo)
		err := storage.AccountAdd(cli.Label, password, account)
		checkErr(err)

		fmt.Printf("'%s' added!\n", cli.Label)

	case "gen":
		password := term.TermGetPassword("Password: ")
		code, nextCode, timeLeft, err := storage.AccountGenCode(cli.Label, password)
		checkErr(err)

		if cli.IsRaw {
			fmt.Print(code)
			return
		}

		clipboard.WriteAll(code)
		fmt.Println()
		fmt.Printf("%s (%ds) (copied!)\n", code, timeLeft)
		fmt.Printf("%s (next)\n", nextCode)

	case "qr":
		password := term.TermGetPassword("Password: ")
		uri, err := storage.AccountGenURI(cli.Label, password)
		checkErr(err)

		term.TermAltScreen()
		term.TermGenQR(uri, cli.IsBig)
		fmt.Print("\nPress [Any]... ")
		term.TermKeyPress()
		term.TermExitAltScreen()

	case "list":
		list := storage.AccountList()
		if len(list) == 0 {
			fmt.Println("No accounts found")
			return
		}
		for _, label := range list {
			fmt.Printf("- %s\n", label)
		}

	case "del":
		password := term.TermGetPassword("Password: ")
		err := storage.AccountDelete(cli.Label, password)
		checkErr(err)

		fmt.Printf("'%s' removed!\n", cli.Label)

	case "rename":
		password := term.TermGetPassword("Password: ")
		err := storage.AccountRename(cli.Label, cli.NewLabel, password)
		checkErr(err)

		fmt.Printf("Renamed: '%s' -> '%s'\n", cli.Label, cli.NewLabel)

	case "import":
		password := term.TermGetPassword("Password: ")
		list, err := storage.AccountImport(cli.Path, password)
		checkErr(err)

		fmt.Printf("Imported accounts: \n")
		for _, label := range list {
			fmt.Printf("- %s\n", label)
		}

	case "export":
		password := term.TermGetPassword("Password: ")
		err := storage.AccountExport(cli.Path, password)
		checkErr(err)

		fmt.Printf("Export to '%s'\n", cli.Path)

	case "completion":
		switch cli.Path {
		case "bash":
			fmt.Print(assets.BashComplete)
		case "zsh":
			fmt.Print(assets.ZshComplete)
		case "fish":
			fmt.Print(assets.FishComplete)
		}

	case "help":
		help()

	default:
		fmt.Println("Usage: mfa --help")
	}

}
