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
	"github.com/ambientsound/visp/style"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type ChangeType int

const (
	ChangeNone               ChangeType = iota // noop
	ChangeList                                 // some list has changes in row data, title, or other
	ChangeOption                               // a setting has been changed
	ChangePlayerStateInvalid                   // player state is no longer valid due to a server command
	ChangeDevice                               // playback device changed
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

	// History returns a list with all tracks played back during the current session.
	History() list.List

	// Return the global multibar instance.
	Multibar() *multibar.Multibar

	// Library returns a list of entry points to the Spotify library.
	Library() *spotify_library.List

	// List returns the active list.
	List() list.List

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

	// UI returns the global UI object.
	UI() UI
}
