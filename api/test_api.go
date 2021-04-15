package api

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/visp/clipboard"
	"github.com/ambientsound/visp/db"
	"github.com/ambientsound/visp/input/keys"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/message"
	"github.com/ambientsound/visp/multibar"
	"github.com/ambientsound/visp/player"
	"github.com/ambientsound/visp/song"
	"github.com/ambientsound/visp/songlist"
	"github.com/ambientsound/visp/spotify/library"
	"github.com/ambientsound/visp/spotify/tracklist"
	"github.com/ambientsound/visp/style"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type testAPI struct {
	messages   chan message.Message
	list       list.List
	song       *song.Song
	songlist   songlist.Songlist
	clipboards *clipboard.List
	tracklist  *spotify_tracklist.List
}

func createTestSong() *song.Song {
	s := song.New()
	s.SetTags(mpd.Attrs{
		"artist": "foo",
		"title":  "bar",
	})
	return s
}

func NewTestAPI() API {
	return &testAPI{
		clipboards: clipboard.New(),
		list:       list.New(),
		messages:   make(chan message.Message, 1024),
		song:       createTestSong(),
		songlist:   songlist.New(),
		tracklist:  spotify_tracklist.NewFromTracks([]spotify.FullTrack{}),
	}
}

func (api *testAPI) Authenticate(token *oauth2.Token) error {
	return nil
}

func (api *testAPI) Clipboards() *clipboard.List {
	return api.clipboards
}

func (api *testAPI) Db() *db.List {
	return nil // FIXME
}

func (api *testAPI) Exec(cmd string) error {
	panic("not implemented")
}

func (api *testAPI) Multibar() *multibar.Multibar {
	panic("not implemented")
}

func (api *testAPI) List() list.List {
	return api.list
}

func (api *testAPI) Library() *spotify_library.List {
	return nil // FIXME
}

func (api *testAPI) ListChanged() {
	// FIXME
}

func (api *testAPI) Message(fmt string, a ...interface{}) {
	api.messages <- message.Format(fmt, a...)
}

func (api *testAPI) MpdClient() *mpd.Client {
	return nil // FIXME
}

func (api *testAPI) OptionChanged(key string) {
	// FIXME
}

func (api *testAPI) Options() Options {
	return viper.GetViper()
}

func (api *testAPI) PlayerStatus() player.State {
	return player.State{}
}

func (api *testAPI) Queue() *songlist.Queue {
	return nil // FIXME
}

func (api *testAPI) Quit() {
	return // FIXME
}

func (api *testAPI) Sequencer() *keys.Sequencer {
	return nil // FIXME
}

func (api *testAPI) SetList(lst list.List) {
	api.list = lst
}

func (api *testAPI) Spotify() (*spotify.Client, error) {
	return nil, fmt.Errorf("no spotify")
}

func (api *testAPI) Song() *song.Song {
	return api.song
}

func (api *testAPI) Songlist() songlist.Songlist {
	return api.songlist
}

func (api *testAPI) Songlists() []songlist.Songlist {
	return nil // FIXME
}

func (api *testAPI) Styles() style.Stylesheet {
	return nil // FIXME
}

func (api *testAPI) Tracklist() *spotify_tracklist.List {
	return api.tracklist
}

func (api *testAPI) UI() UI {
	return nil // FIXME
}
