package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/ambientsound/visp/log"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/google/uuid"
)

// Duplicate makes a copy of a tracklist.
type Duplicate struct {
	command
	api       api.API
	name      string
	tracklist *spotify_tracklist.List
}

// NewDuplicate returns Duplicate.
func NewDuplicate(api api.API) Command {
	return &Duplicate{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Duplicate) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteEmpty()

	cmd.tracklist = cmd.api.Tracklist()
	if cmd.tracklist == nil {
		return fmt.Errorf("only track lists can be duplicated")
	}

	switch tok {
	case lexer.TokenEnd:
		// No parameters; save original list back to itself
		cmd.name = "Copy of " + cmd.tracklist.Name()
	case lexer.TokenIdentifier:
		// New name; write this list to a new copy
		cmd.name = lit
	default:
		return fmt.Errorf("unexpected '%s', expected name of new playlist", lit)
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Duplicate) Exec() error {
	tracklist := cmd.tracklist.Copy()
	tracklist.SetID(uuid.New().String())
	tracklist.SetName(cmd.name)
	cmd.api.Db().Cache(tracklist)

	log.Infof("Created '%s' with %d tracks", tracklist.Name(), tracklist.Len())

	return nil
}
