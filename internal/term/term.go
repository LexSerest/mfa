package term

import (
	"fmt"
	"os"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"golang.org/x/term"
)

func TermGetPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
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
