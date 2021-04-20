package commands

import (
	"fmt"
	"strconv"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Rename saves a local tracklist to Spotify.
type Rename struct {
	command
	api  api.API
	name string
}

// NewRename returns Rename.
func NewRename(api api.API) Command {
	return &Rename{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Rename) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabComplete(lit, []string{strconv.Quote(cmd.api.List().Name())})

	switch tok {
	case lexer.TokenIdentifier:
		cmd.name = lit
	default:
		return fmt.Errorf("unexpected '%s', expected new name", lit)
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Rename) Exec() error {
	cmd.api.List().SetName(cmd.name)
	cmd.api.Changed(api.ChangeList, cmd.api.List())
	return nil
}
