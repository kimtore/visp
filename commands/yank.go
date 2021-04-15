package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/log"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
)

// Yank copies tracks from the songlist into the clipboard.
type Yank struct {
	command
	api  api.API
	list *spotify_tracklist.List
}

// NewYank returns Yank.
func NewYank(api api.API) Command {
	return &Yank{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Yank) Parse() error {
	cmd.list = cmd.api.Tracklist()
	if cmd.list == nil {
		return fmt.Errorf("`yank` only works in tracklists")
	}
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Yank) Exec() error {
	selection := cmd.list.Selection()

	if selection.Len() == 0 {
		return fmt.Errorf("no tracks selected")
	}

	selection.SetVisibleColumns(cmd.list.ColumnNames())

	cmd.api.Clipboards().Insert(&selection)
	log.Infof("%d songs stored in %s", selection.Len(), selection.Name())

	cmd.list.ClearSelection()

	return nil
}
