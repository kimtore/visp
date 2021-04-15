package commands

import (
	"fmt"

	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/zmb3/spotify"

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
	tracklist *spotify_tracklist.List
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
	cmd.tracklist = cmd.api.Tracklist()

	if err != nil {
		return err
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
		return cmd.client.Play()
	}
}

// playCursor plays the song under the cursor.
func (cmd *Play) playCursor() error {
	if cmd.tracklist == nil {
		return fmt.Errorf("cannot play cursor when not in a track list")
	}

	// Get the song under the cursor.
	track := cmd.tracklist.CursorTrack()
	if track == nil {
		return fmt.Errorf("cannot play: no track under cursor")
	}

	// Play the correct song.
	return cmd.client.PlayOpt(&spotify.PlayOptions{
		URIs: []spotify.URI{
			track.URI,
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

	cmd.tracklist.ClearSelection()

	uris := make([]spotify.URI, selection.Len())
	for i, track := range selection.Tracks() {
		uris[i] = track.URI
	}

	// TODO: queue is unsupported by the Spotify Web API
	// https://github.com/spotify/web-api/issues/462

	// Start playing all the URIs
	return cmd.client.PlayOpt(&spotify.PlayOptions{
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
