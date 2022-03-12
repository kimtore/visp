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
	current   bool
	cursor    bool
	selection bool
	// mode of operation
	add    bool
	remove bool
	toggle bool
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
	case "toggle":
		cmd.toggle = true
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
	case "current":
		cmd.current = true
	case "cursor":
		cmd.cursor = true
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

	tracklist := cmd.api.List()

	if !cmd.current && tracklist == nil {
		return fmt.Errorf("liking tracks by cursor or selection needs an active tracklist")
	}

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	ids := make([]spotify.ID, 0)

	switch {
	case cmd.cursor:
		track := tracklist.CursorRow()
		if track != nil {
			ids = append(ids, spotify.ID(track.ID()))
		}
	case cmd.selection:
		tracks := tracklist.Selection().All()
		for _, track := range tracks {
			ids = append(ids, spotify.ID(track.ID()))
		}
	case cmd.current:
		id := cmd.api.PlayerStatus().TrackRow.ID()
		if len(id) == 0 {
			return fmt.Errorf("no track is playing right now")
		}
		ids = append(ids, spotify.ID(id))
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	additions := make([]spotify.ID, 0, len(ids))
	removals := make([]spotify.ID, 0, len(ids))

	switch {
	case cmd.add:
		additions = append(additions, ids...)
	case cmd.remove:
		removals = append(removals, ids...)
	case cmd.toggle:
		liked, err := client.UserHasTracks(ids...)
		if err != nil {
			return err
		}
		for i, trackLiked := range liked {
			if trackLiked {
				removals = append(removals, ids[i])
			} else {
				additions = append(additions, ids[i])
			}
		}
	}

	if len(additions) > 0 {
		err = client.AddTracksToLibrary(additions...)
		if err != nil {
			return err
		}
		log.Infof("%d track(s) added to Liked tracks", len(additions))
	}

	if len(removals) > 0 {
		err = client.RemoveTracksFromLibrary(removals...)
		if err != nil {
			return err
		}
		log.Infof("%d track(s) removed from Liked tracks", len(removals))
	}

	if cmd.selection {
		tracklist.ClearSelection()
	}

	return nil
}

func (cmd *Like) setTabCompleteSelection(lit string) {
	cmd.setTabComplete(lit, []string{
		"current",
		"cursor",
		"selection",
	})
}

func (cmd *Like) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"add",
		"remove",
		"toggle",
	})
}
