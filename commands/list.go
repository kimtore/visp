package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/spotify/aggregator"
	"github.com/ambientsound/visp/spotify/devices"
	"github.com/ambientsound/visp/spotify/library"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// List navigates and manipulates songlists.
type List struct {
	command
	api       api.API
	client    *spotify.Client
	absolute  int
	duplicate bool
	goto_     bool
	open      bool
	new       bool
	relative  int
	close     bool
	last      bool
	name      string
}

func NewList(api api.API) Command {
	return &List{
		api:      api,
		absolute: -1,
	}
}

func (cmd *List) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
		switch lit {
		case "duplicate":
			cmd.duplicate = true
			cmd.name = cmd.api.List().Name()
		case "close":
			cmd.close = true
		case "up", "prev", "previous":
			cmd.relative = -1
		case "down", "next":
			cmd.relative = 1
		case "home":
			cmd.absolute = 0
		case "end":
			cmd.absolute = cmd.api.Db().Len() - 1
		case "goto":
			cmd.goto_ = true
		case "open":
			cmd.open = true
		case "last":
			cmd.last = true
		case "new":
			cmd.new = true
		default:
			i, err := strconv.Atoi(lit)
			if err != nil {
				return fmt.Errorf("cannot navigate lists: position '%s' is not recognized, and is not a number", lit)
			}
			cmd.absolute = i - 1
		}
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	if cmd.goto_ || cmd.new {
		for tok != lexer.TokenEnd {
			tok, lit = cmd.Scan()
			cmd.name += lit
		}

		cmd.name = strings.TrimSpace(cmd.name)
		cmd.Unscan()

		if cmd.goto_ {
			cmd.setTabComplete(cmd.name, cmd.api.Db().Keys())
		}
	} else {
		cmd.setTabCompleteEmpty()
	}

	if cmd.new && len(cmd.name) == 0 {
		cmd.name = cmd.generateName()
	}

	return cmd.ParseEnd()
}

func (cmd *List) Exec() error {
	switch {
	case cmd.goto_:
		return cmd.Goto(cmd.name)

	case cmd.last:
		cmd.api.SetList(cmd.api.Db().Last())
		return nil

	case cmd.open:
		row := cmd.api.List().CursorRow()
		if row == nil {
			return fmt.Errorf("no playlist selected")
		}
		return cmd.Goto(row.ID())

	case cmd.relative != 0:
		cmd.api.Db().MoveCursor(cmd.relative)
		cmd.api.SetList(cmd.api.Db().Current())

	case cmd.absolute >= 0:
		cmd.api.Db().SetCursor(cmd.absolute)
		cmd.api.SetList(cmd.api.Db().Current())

	case cmd.duplicate:
		return cmd.Duplicate()

	case cmd.new:
		return cmd.New()

	case cmd.close:
		db := cmd.api.Db()
		cur := db.Current()
		err := db.Remove(db.Cursor())
		if err != nil {
			return err
		}
		if cur != nil {
			log.Infof("Closed '%s'", cur.Name())
		}
		if db.Len() == 0 {
			db.Cache(log.List(log.InfoLevel))
		}
		db.SetCursor(db.Cursor())
		cmd.api.SetList(db.Current())
	}

	return nil
}

// Goto loads an external list and applies default columns and sorting.
// Local, cached versions are tried first.
func (cmd *List) Goto(id string) error {
	var err error
	var lst list.List

	// Set Spotify object request limit. Ignore user-defined max limit here,
	// because big queries will always be faster and consume less bandwidth,
	// when requesting all the data.
	const limit = 50

	// Try a cached version of a named list
	lst = cmd.api.Db().List(cmd.name)
	if lst != nil {
		cmd.api.SetList(lst)
		return nil
	}

	// Other named lists need Spotify access
	cmd.client, err = cmd.api.Spotify()
	if err != nil {
		return err
	}

	t := time.Now()
	switch id {
	case spotify_library.MyPlaylists:
		lst, err = spotify_aggregator.MyPrivatePlaylists(*cmd.client, limit)
	case spotify_library.FeaturedPlaylists:
		lst, err = spotify_aggregator.FeaturedPlaylists(*cmd.client, limit)
	case spotify_library.MyTracks:
		lst, err = spotify_aggregator.MyTracks(*cmd.client, limit)
	case spotify_library.TopTracks:
		lst, err = spotify_aggregator.TopTracks(*cmd.client, limit)
	case spotify_library.NewReleases:
		lst, err = spotify_aggregator.NewReleases(*cmd.client)
	case spotify_library.MyAlbums:
		lst, err = spotify_aggregator.MyAlbums(*cmd.client)
	case spotify_library.Devices:
		lst, err = spotify_devices.New(*cmd.client)
	default:
		lst, err = spotify_aggregator.ListWithID(*cmd.client, id, limit)
		if err != nil {
			break
		}
	}
	dur := time.Since(t)

	if err != nil {
		return err
	}

	log.Debugf("Retrieved %s with %d items in %s", id, lst.Len(), dur.String())
	log.Infof("Loaded %s.", lst.Name())

	// Reset cursor
	lst.SetCursor(0)

	cmd.api.SetList(lst)

	return nil
}

func (cmd *List) Duplicate() error {
	tracklist := cmd.api.List().Copy()
	tracklist.SetName(cmd.name)
	tracklist.SetID(uuid.New().String())
	tracklist.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	cmd.api.SetList(tracklist)

	log.Infof("Created temporary playlist '%s' with %d tracks", tracklist.Name(), tracklist.Len())

	return nil
}

var (
	totalNewCreated = 0
)

// Exec implements Command.
func (cmd *List) New() error {
	tracklist := spotify_tracklist.NewFromTracks([]spotify.FullTrack{})
	tracklist.SetName(cmd.name)
	tracklist.SetID(uuid.New().String())
	tracklist.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	cmd.api.SetList(tracklist)

	log.Infof("Created temporary playlist '%s'", tracklist.Name())

	return nil
}

// Generates a new playlist name.
func (cmd *List) generateName() string {
	totalNewCreated++
	return fmt.Sprintf("New playlist %d", totalNewCreated)
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *List) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"close",
		"down",
		"duplicate",
		"end",
		"goto",
		"home",
		"last",
		"new",
		"next",
		"prev",
		"previous",
		"up",
	})
}
