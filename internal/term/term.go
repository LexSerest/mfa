package term

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"golang.org/x/term"
)

func TermGetPassword(prompt string) string {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Print(prompt)
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		return string(passwordBytes)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from pipeline: %v\n", err)
		os.Exit(1)
	}

	return ""
}

func TermAltScreen() {
	fmt.Print("\033[?1049h\033[H")
}

func TermExitAltScreen() {
	fmt.Print("\033[?1049l")
}

func TermKeyPress() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err == nil {
		buf := make([]byte, 1)
		_, _ = os.Stdin.Read(buf)
		_ = term.Restore(fd, oldState)
	} else {
		_, _ = os.Stdin.Read(make([]byte, 1))
	}
}

func TermGenQR(uri string, isBig bool) {
	if isBig {
		qrterminal.Generate(uri, qrterminal.L, os.Stdout)
	} else {
		qrterminal.GenerateHalfBlock(uri, qrterminal.L, os.Stdout)
	}
}
