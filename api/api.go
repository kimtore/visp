// Package api provides data model interfaces.
package api

import (
	"github.com/ambientsound/visp/clipboard"
	"github.com/ambientsound/visp/db"
	"github.com/ambientsound/visp/input/keys"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/multibar"
	"github.com/ambientsound/visp/player"
	"github.com/ambientsound/visp/spotify/library"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/ambientsound/visp/style"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type ChangeType int

const (
	ChangeNone ChangeType = iota // noop
	ChangeList
	ChangeOption
)

// API defines a set of commands that should be available to commands run
// through the command-line interface.
type API interface {
	// Authenticate sets an OAuth2 token that should be used for Spotify calls.
	Authenticate(token *oauth2.Token) error

	// Changed notifies the program that some internal state has changed.
	Changed(typ ChangeType, data interface{})

	// Clipboards is a list of clipboards.
	Clipboards() *clipboard.List

	// Db returns the database of lists.
	Db() *db.List

	// Exec executes a command through the command-line interface.
	Exec(string) error

	// Return the global multibar instance.
	Multibar() *multibar.Multibar

	// Library returns a list of entry points to the Spotify library.
	Library() *spotify_library.List

	// List returns the active list.
	List() list.List

	// Message sends a message to the user through the statusbar.
	Message(string, ...interface{})

	// Options returns PMS' global options.
	Options() Options

	// PlayerStatus returns the current MPD player status.
	PlayerStatus() player.State

	// Quit shuts down PMS.
	Quit()

	// Sequencer returns a pointer to the key sequencer that receives key events.
	Sequencer() *keys.Sequencer

	// SetList sets the active list.
	SetList(list.List)

	// Spotify returns a Spotify client.
	Spotify() (*spotify.Client, error)

	// Styles returns the current stylesheet.
	Styles() style.Stylesheet

	// Tracklist returns the visible track list, if any.
	// Will be nil if the active widget shows a different kind of list.
	Tracklist() *spotify_tracklist.List

	// UI returns the global UI object.
	UI() UI
}
