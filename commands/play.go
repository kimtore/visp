package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/zmb3/spotify/v2"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Play plays songs in the MPD playlist.
type Play struct {
	command
	api       api.API
	cursor    bool
	selection bool
	client    *spotify.Client
	tracklist list.List
}

// NewPlay returns Play.
func NewPlay(api api.API) Command {
	return &Play{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Play) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenEnd:
		// No parameters; just send 'play' command to MPD
		return nil
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	switch lit {
	// Play song under cursor
	case "cursor":
		cmd.cursor = true
	// Play selected songs
	case "selection":
		cmd.selection = true
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Play) Exec() error {
	var err error

	cmd.client, err = cmd.api.Spotify()
	cmd.tracklist = cmd.api.List()

	if err != nil {
		return err
	}

	if cmd.tracklist == nil || cmd.tracklist.Len() == 0 {
		return fmt.Errorf("no tracks in active list")
	}

	switch {
	case cmd.cursor:
		// Play song under cursor.
		return cmd.playCursor()
	case cmd.selection:
		// Play selected songs.
		return cmd.playSelection()
	default:
		// If a selection is not given, start playing with default parameters.
		return cmd.client.Play(context.TODO())
	}
}

// If no device is active, use the `device` option to automatically start playing
// on a preferred device.
func (cmd *Play) deviceID() (*spotify.ID, error) {
	if len(cmd.api.PlayerStatus().Device.ID) > 0 {
		return nil, nil
	}

	deviceName := options.GetString(options.Device)
	if len(deviceName) == 0 {
		return nil, nil
	}

	devices, err := cmd.client.PlayerDevices(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to determine device ID for playback: %w", err)
	}

	for _, device := range devices {
		if device.Name == deviceName {
			return &device.ID, nil
		}
	}

	// Fallback to lower-case matching if no exact match is found
	for _, device := range devices {
		if strings.ToLower(device.Name) == strings.ToLower(deviceName) {
			return &device.ID, nil
		}
	}

	return nil, fmt.Errorf("no active playback device and no default set; try `list goto devices` or `set device=<device_name>`")
}

// playCursor plays the song under the cursor, and also adds the rest of the list to the play context.
// Playback starts from the song beneath the cursor.
func (cmd *Play) playCursor() error {
	row := cmd.tracklist.CursorRow()

	if row.Kind() != list.DataTypeTrack {
		return fmt.Errorf("cannot play: %w", ErrMsgDataType(row.Kind(), list.DataTypeTrack))
	}

	// Get a device ID for playback.
	deviceID, err := cmd.deviceID()
	if err != nil {
		return err
	}

	// Figure out the playback context.
	// If the local list is a remote playlist, AND it is in sync with the server,
	// we can use the playlist URI as playback context in order to have the queue
	// display our played tracks in the desktop app.
	//
	// If the local list has changes, or it is not a remote playlist, we compose
	// a list of all tracks in the list, and use this ad-hoc list as a playback context.
	//
	// Unfortunately, Spotify has a limit on how many tracks can be played ad-hoc.
	// This seems to be based on the size of the HTTP request, and not the API itself.
	// Any request too large will return with 'HTTP 413 Request Entity Too Large'.
	// Simple tests placed this number at 784 tracks.
	// Here, we set the limit to 750 tracks to be sure.
	uri := cmd.tracklist.URI()
	uris := make([]spotify.URI, 0, cmd.tracklist.Len())
	if uri == nil || cmd.tracklist.HasLocalChanges() {
		const limit = 750
		uri = nil
		tracks := cmd.tracklist.All()
		if len(tracks) > limit {
			log.Infof("Note: tracklist contains %d tracks, but only %d tracks will be added to avoid errors", len(tracks), limit)
			tracks = tracks[:limit]
		}
		for _, tr := range tracks {
			uris = append(uris, tr.URI())
		}
		log.Infof("Starting playback of %d tracks starting with '%s'", len(uris), row.Get("title"))
	} else {
		uris = nil
		log.Infof("Starting playback of playlist '%s', starting with '%s'", cmd.tracklist.Name(), row.Get("title"))
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)
	defer cmd.api.Changed(api.ChangeDevice, nil)

	// Start playing with correct parameters.
	trackuri := spotify.URI("spotify:track:" + row.ID())
	return cmd.client.PlayOpt(context.TODO(), &spotify.PlayOptions{
		DeviceID:        deviceID,
		URIs:            uris,
		PlaybackContext: uri,
		PlaybackOffset: &spotify.PlaybackOffset{
			URI: trackuri,
		},
	})
}

// playSelection plays the currently selected songs.
func (cmd *Play) playSelection() error {

	if cmd.tracklist == nil {
		return fmt.Errorf("cannot play cursor when not in a track list")
	}

	selection := cmd.tracklist.Selection()
	if selection.Len() == 0 {
		return fmt.Errorf("cannot play: no selection")
	}

	// Selection is cursor
	row := cmd.tracklist.CursorRow()
	if selection.Len() == 1 && row != nil && selection.Row(0).ID() == row.ID() {
		return cmd.playCursor()
	}

	cmd.tracklist.ClearSelection()

	uris := make([]spotify.URI, selection.Len())
	for i, track := range selection.All() {
		uris[i] = track.URI()
	}

	// TODO: queue is unsupported by the Spotify Web API
	// https://github.com/spotify/web-api/issues/462

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	// Start playing all the URIs
	return cmd.client.PlayOpt(context.TODO(), &spotify.PlayOptions{
		URIs: uris,
	})
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Play) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"cursor",
		"selection",
	})
}
