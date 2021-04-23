package commands

import (
	"fmt"

	"github.com/ambientsound/visp/log"
	"github.com/zmb3/spotify"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Like plays songs in the MPD playlist.
type Like struct {
	command
	api api.API
	// selector
	cursor    bool
	playing   bool
	selection bool
	// mode of operation
	add    bool
	remove bool
}

// NewLike returns Like.
func NewLike(api api.API) Command {
	return &Like{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Like) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case "add":
		cmd.add = true
	case "remove":
		cmd.remove = true
	default:
		return fmt.Errorf("unexpected '%s', expected 'add' or 'remove'", lit)
	}

	cmd.setTabCompleteEmpty()

	tok, lit = cmd.Scan()
	if tok != lexer.TokenWhitespace {
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	tok, lit = cmd.Scan()
	cmd.setTabCompleteSelection(lit)

	switch tok {
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case "cursor":
		cmd.cursor = true
	case "playing":
		cmd.playing = true
	case "selection":
		cmd.selection = true
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Like) Exec() error {
	var err error

	tracklist := cmd.api.Tracklist()

	if !cmd.playing && tracklist == nil {
		return fmt.Errorf("liking tracks by cursor or selection needs an active tracklist")
	}

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	ids := make([]spotify.ID, 0)
	switch {
	case cmd.cursor:
		track := tracklist.CursorTrack()
		if track != nil {
			ids = append(ids, track.ID)
		}
	case cmd.selection:
		tracks := tracklist.Selection().Tracks()
		for _, track := range tracks {
			ids = append(ids, track.ID)
		}
	case cmd.playing:
		id := cmd.api.PlayerStatus().TrackRow.ID()
		if len(id) == 0 {
			return fmt.Errorf("no track is playing right now")
		}
		ids = append(ids, spotify.ID(id))
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	if cmd.add {
		err = client.AddTracksToLibrary(ids...)
	} else if cmd.remove {
		err = client.RemoveTracksFromLibrary(ids...)
	}

	if err != nil {
		return err
	}

	if cmd.add {
		log.Infof("%d track(s) added to Liked tracks", len(ids))
	} else if cmd.remove {
		log.Infof("%d track(s) removed from Liked tracks", len(ids))
	}

	if cmd.selection {
		tracklist.ClearSelection()
	}

	return nil
}

func (cmd *Like) setTabCompleteSelection(lit string) {
	cmd.setTabComplete(lit, []string{
		"cursor",
		"playing",
		"selection",
	})
}

func (cmd *Like) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"add",
		"remove",
	})
}