package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
)

// Yank copies tracks from the songlist into the clipboard.
type Yank struct {
	command
	api     api.API
	current bool
	list    list.List
}

// NewYank returns Yank.
func NewYank(api api.API) Command {
	return &Yank{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Yank) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
		if lit == "current" {
			cmd.current = true
		} else {
			cmd.Unscan()
		}
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Yank) Exec() error {
	switch {
	case cmd.current == true:
		row := cmd.api.PlayerStatus().TrackRow
		if len(row.ID()) == 0 {
			return fmt.Errorf("no track currently playing")
		}
		cmd.list = list.New()
		cmd.list.Add(row)
		cmd.list.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	default:
		tracklist := cmd.api.Tracklist()
		if tracklist == nil {
			return fmt.Errorf("`yank` only works in tracklists")
		}
		cmd.list = tracklist.Selection()

		if cmd.list.Len() == 0 {
			return fmt.Errorf("no tracks selected")
		}

		cmd.list.SetVisibleColumns(tracklist.ColumnNames())
		tracklist.ClearSelection()
	}

	cmd.api.Clipboards().Insert(cmd.list)
	log.Infof("%d songs stored in %s", cmd.list.Len(), cmd.list.Name())

	return nil
}

func (cmd *Yank) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"current",
	})
}
