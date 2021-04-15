package commands

import (
	"fmt"

	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

var (
	totalNewCreated = 0
)

// NewCmd creates a new track list.
// The Cmd suffix is there because the name is taken.
type NewCmd struct {
	command
	api  api.API
	name string
}

// NewNew returns NewCmd.
func NewNew(api api.API) Command {
	return &NewCmd{
		api: api,
	}
}

// Parse implements Command.
func (cmd *NewCmd) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteEmpty()

	switch tok {
	case lexer.TokenEnd:
		// No parameters; create a new list with a generic name
		cmd.name = cmd.generateName()
	case lexer.TokenIdentifier:
		cmd.name = lit
	default:
		return fmt.Errorf("unexpected '%s', expected name of playlist", lit)
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *NewCmd) Exec() error {
	tracklist := spotify_tracklist.NewFromTracks([]spotify.FullTrack{})
	tracklist.SetName(cmd.name)
	tracklist.SetID(uuid.New().String())
	tracklist.SetVisibleColumns(options.GetList(options.Columns))

	cmd.api.SetList(tracklist)

	return nil
}

// Generates a new playlist name.
func (cmd *NewCmd) generateName() string {
	totalNewCreated++
	return fmt.Sprintf("New playlist %d", totalNewCreated)
}
