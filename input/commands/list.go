package commands

import (
	"fmt"
	"strconv"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// List navigates songlists.
type List struct {
	ui       *widgets.UI
	relative int
	absolute int
}

func NewList(ui *widgets.UI) *List {
	return &List{ui: ui}
}

func (cmd *List) Reset() {
	cmd.relative = 0
	cmd.absolute = -1
}

func (cmd *List) Execute(t lexer.Token) error {
	var err error

	s := t.String()

	switch t.Class {

	case lexer.TokenIdentifier:
		switch s {
		case "up", "prev", "previous":
			cmd.relative = -1
		case "down", "next":
			cmd.relative = 1
		case "home":
			cmd.absolute = 0
		case "end":
			cmd.absolute = cmd.ui.SonglistsLen() - 1
		default:
			i, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("Cannot navigate lists: position '%s' is not recognized, and is not a number", s)
			}
			switch {
			case cmd.relative != 0 || cmd.absolute != -1:
				return fmt.Errorf("Only one number allowed when setting list position")
			case cmd.relative != 0:
				cmd.relative *= i
			default:
				cmd.absolute = i - 1
			}
		}

	case lexer.TokenEnd:
		switch {
		case cmd.relative != 0:
			index := cmd.ui.SonglistIndex() + cmd.relative
			if !cmd.ui.ValidSonglistIndex(index) {
				len := cmd.ui.SonglistsLen()
				index = (index + len) % len
			}
			console.Log("Setting tab to relative %d", index)
			err = cmd.ui.SetSonglistIndex(index)
		case cmd.absolute >= 0:
			console.Log("Setting tab to absolute %d", cmd.absolute)
			err = cmd.ui.SetSonglistIndex(cmd.absolute)
		default:
			return fmt.Errorf("Unexpected END, expected position. Try one of: next prev <number>")
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
