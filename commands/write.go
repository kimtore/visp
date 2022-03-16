package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ambientsound/visp/log"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/zmb3/spotify/v2"

	"github.com/ambientsound/visp/api"
)

// Write saves a local tracklist to Spotify.
type Write struct {
	command
	api           api.API
	name          string
	new           bool
	public        bool
	collaborative bool // TODO: implement this
}

// NewWrite returns Write.
func NewWrite(api api.API) Command {
	return &Write{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Write) Parse() error {
	lit := cmd.ScanRemainderAsIdentifier()

	cmd.setTabComplete(lit, []string{strconv.Quote(cmd.api.List().Name())})

	if len(lit) > 0 {
		cmd.name = lit
	}

	// TODO: private/public?

	return nil
}

// Exec implements Command.
func (cmd *Write) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	tracklist := cmd.api.List()

	// Copy tracklist, assign new name, and save that one
	if len(cmd.name) > 0 {
		tracklist = tracklist.Copy()
		tracklist.SetName(cmd.name)
		cmd.api.Db().Cache(tracklist)
	}

	row := cmd.api.Db().RowByID(tracklist.ID())
	if row == nil {
		return fmt.Errorf("internal error: tracklist not cached in database")
	}

	user, err := client.CurrentUser(context.TODO())
	if err != nil {
		return err
	}

	ids := make([]spotify.ID, 0, tracklist.Len())
	for _, track := range tracklist.All() {
		ids = append(ids, spotify.ID(track.ID()))
	}

	if !tracklist.HasRemote() {
		remotelist, err := client.CreatePlaylistForUser(context.TODO(), user.ID, tracklist.Name(), "", cmd.public, cmd.collaborative)
		if err != nil {
			return fmt.Errorf("create remote playlist: %w", err)
		}

		tracklist.SetID(remotelist.ID.String())
		tracklist.SetName(remotelist.Name)
		tracklist.SetRemote(true)
		row.SetID(tracklist.ID())

		// Re-index original list in database if working on the old copy

		cmd.api.SetList(tracklist)

		snapshot, err := spotify_tracklist.AddTracksToPlaylist(client, remotelist.ID, ids)
		if err != nil {
			return fmt.Errorf("add tracks to remote playlist: %w", err)
		}
		_ = snapshot // todo: add and use this?

		log.Infof("Created playlist '%s' with %d tracks", remotelist.Name, len(ids))

	} else {

		id := spotify.ID(tracklist.ID())
		err := client.ChangePlaylistName(context.TODO(), id, tracklist.Name())
		if err != nil {
			return fmt.Errorf("change remote playlist name: %w", err)
		}

		err = spotify_tracklist.ReplacePlaylistTracks(client, id, ids)
		if err != nil {
			return fmt.Errorf("write new track list to to remote playlist: %w", err)
		}

		log.Infof("Wrote changes to remote playlist '%s' with %d tracks", tracklist.Name(), len(ids))
	}

	tracklist.SetSyncedToRemote()

	return nil
}
