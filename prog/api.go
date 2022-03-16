package prog

import (
	"context"
	"fmt"
	"time"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/clipboard"
	"github.com/ambientsound/visp/db"
	"github.com/ambientsound/visp/input/keys"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/multibar"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/player"
	"github.com/ambientsound/visp/spotify/library"
	"github.com/ambientsound/visp/spotify/proxyclient"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/ambientsound/visp/style"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (v *Visp) Authenticate(token *oauth2.Token) error {
	log.Infof("Configured Spotify access token, expires at %s", token.Expiry.Format(time.RFC1123))

	cli := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	v.client = spotify.New(cli)

	next := spotify_proxyclient.TokenTTL(token)

	v.tokenRefresh = time.After(next)

	v.Changed(api.ChangePlayerStateInvalid, nil)

	err := v.Tokencache.Write(*token)
	if err != nil {
		return fmt.Errorf("write Spotify token to file: %s", err)
	}

	return nil
}

func (v *Visp) Clipboards() *clipboard.List {
	return v.clipboards
}

func (v *Visp) Db() *db.List {
	return v.db
}

func (v *Visp) Exec(command string) error {
	log.Debugf("Run command: %s", command)
	return v.interpreter.Exec(command)
}

func (v *Visp) Library() *spotify_library.List {
	return v.library
}

func (v *Visp) List() list.List {
	return v.list
}

func (v *Visp) Changed(change api.ChangeType, data interface{}) {
	switch change {
	case api.ChangeList:
		lst, ok := data.(list.List)
		if !ok {
			log.Debugf("BUG: list was changed, but is '%T', not 'list.List'", data)
			return
		}
		v.db.Cache(lst)
		v.clipboards.Update(lst)
		// TODO: playlists indexes

	case api.ChangeOption:
		s, ok := data.(string)
		if !ok {
			log.Debugf("BUG: option '#v' changed, but is not of type 'string'", data)
			return
		}
		v.optionChanged(s)

	case api.ChangePlayerStateInvalid:
		v.player.Invalidate()
		v.ticker.Reset(changePlayerStateDelay)

	case api.ChangeDevice:
		v.player.Invalidate()
		v.ticker.Reset(changePlayerStateDelay)
		// TODO: refresh devices window

	}
}

func (v *Visp) PlayerStatus() player.State {
	return *v.player
}

func (v *Visp) Quit() {
	v.quit <- new(interface{})
}

func (v *Visp) Sequencer() *keys.Sequencer {
	return v.sequencer
}

func (v *Visp) Multibar() *multibar.Multibar {
	return v.multibar
}

func (v *Visp) History() list.List {
	if v.history == nil {
		v.history = spotify_tracklist.NewHistory()
	}
	return v.history
}

func (v *Visp) SetList(lst list.List) {
	if lst == nil {
		return
	}
	cur := v.db.Current()
	if cur != nil && cur != lst && cur != v.db && cur != v.clipboards {
		log.Debugf("Setting last used list to '%s'", cur.Name())
		v.db.SetLast(v.db.Current())
	}
	c := v.db.Cache(lst)
	v.db.SetCursor(c)
	v.list = lst
	v.Termui.TableWidget().SetList(lst)
}

func (v *Visp) Spotify() (*spotify.Client, error) {
	if v.client == nil {
		return nil, fmt.Errorf("please authenticate with Spotify at: %s/authorize", options.GetString("spotifyauthserver"))
	}
	token, err := v.client.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh Spotify token: %s", err)
	}
	_ = v.Tokencache.Write(*token)
	return v.client, nil
}

func (v *Visp) Styles() style.Stylesheet {
	return v.stylesheet
}

func (v *Visp) UI() api.UI {
	return v.Termui
}
