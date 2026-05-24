package assets

import (
	_ "embed"
)

//go:embed help.txt
var HelpText string

//go:embed completion.fish
var FishComplete string

//go:embed completion.zsh
var ZshComplete string

//go:embed completion.bash
var BashComplete string
