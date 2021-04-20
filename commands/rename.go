package commands

import (
	"strconv"

	"github.com/ambientsound/visp/api"
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
	lit := cmd.ScanRemainderAsIdentifier()

	cmd.setTabComplete(lit, []string{strconv.Quote(cmd.api.List().Name())})
	cmd.name = lit

	return nil
}

// Exec implements Command.
func (cmd *Rename) Exec() error {
	cmd.api.List().SetName(cmd.name)
	cmd.api.Changed(api.ChangeList, cmd.api.List())
	return nil
}
