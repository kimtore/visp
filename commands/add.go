package commands

import (
	"context"
	"fmt"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/zmb3/spotify/v2"

	"github.com/ambientsound/visp/api"
)

// Add plays songs in the MPD playlist.
type Add struct {
	command
	api       api.API
	client    *spotify.Client
	tracklist list.List
}

// NewAdd returns Add.
func NewAdd(api api.API) Command {
	return &Add{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Add) Parse() error {
	cmd.tracklist = cmd.api.List()
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Add) Exec() error {
	var err error

	cmd.client, err = cmd.api.Spotify()

	if err != nil {
		return err
	}

	selection := cmd.tracklist.Selection()
	if selection.Len() == 0 {
		return fmt.Errorf("cannot add to queue: no selection")
	}

	// Allow command to deselect tracks in visual selection that were added to the queue.
	// In case of a queue add failure, it is desirable to still select the tracks that failed
	// to be added.
	cmd.tracklist.CommitVisualSelection()
	cmd.tracklist.DisableVisualSelection()

	tracks := selection.All()
	for i, track := range tracks {
		err = ErrMsgDataType(track.Kind(), list.DataTypeTrack)
		if err != nil {
			return err
		}

		err = cmd.client.QueueSong(context.TODO(), spotify.ID(track.ID()))
		if err != nil {
			return err
		}
		log.Infof("'%s - %s' added to queue.", track.Get("artist"), track.Get("title"))
		cmd.tracklist.SetSelected(i, false)
	}

	return nil
}
