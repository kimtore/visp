package commands

import (
	"fmt"

	"github.com/ambientsound/visp/log"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/zmb3/spotify"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Write saves a local tracklist to Spotify.
type Write struct {
	command
	api     api.API
	name    string
	new     bool
	private bool
}

// NewWrite returns Write.
func NewWrite(api api.API) Command {
	return &Write{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Write) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteEmpty()

	switch tok {
	case lexer.TokenEnd:
		// No parameters; save original list back to itself
	case lexer.TokenIdentifier:
		// New name; write this list to a new copy
		cmd.name = lit
	default:
		return fmt.Errorf("unexpected '%s', expected name of playlist", lit)
	}

	// TODO
	cmd.private = true

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Write) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	tracklist := cmd.api.Tracklist()
	if tracklist == nil {
		return fmt.Errorf("only track lists can be saved to Spotify")
	}

	// Copy tracklist, assign new name, and save that one
	if len(cmd.name) > 0 {
		tracklist = tracklist.Copy()
		tracklist.SetName(cmd.name)
	}

	row := cmd.api.Db().RowByID(tracklist.ID())
	if row == nil {
		return fmt.Errorf("internal error: can't find tracklist in local database")
	}

	user, err := client.CurrentUser()
	if err != nil {
		return err
	}

	if !tracklist.HasRemote() {
		remotelist, err := client.CreatePlaylistForUser(user.ID, tracklist.Name(), "", cmd.private)
		if err != nil {
			return fmt.Errorf("create remote playlist: %w", err)
		}

		tracklist.SetID(remotelist.ID.String())
		tracklist.SetName(remotelist.Name)
		tracklist.SetRemote(true)

		// Re-index original list in database if working on the old copy
		if len(cmd.name) == 0 {
			row.SetID(tracklist.ID())
		}

		ids := make([]spotify.ID, 0, tracklist.Len())
		for _, track := range tracklist.Tracks() {
			ids = append(ids, track.ID)
		}

		snapshot, err := spotify_tracklist.AddTracksToPlaylist(client, remotelist.ID, ids)
		if err != nil {
			return fmt.Errorf("add tracks to remote playlist: %w", err)
		}
		_ = snapshot // todo: add and use this?

		log.Infof("Created playlist '%s' with %d tracks", remotelist.Name, len(ids))

		tracklist.SetSyncedToRemote()

	} else {

		return fmt.Errorf("writing changes to existing lists is unimplemented")
	}

	return nil
}
