package prog

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/clipboard"
	"github.com/ambientsound/visp/commands"
	"github.com/ambientsound/visp/db"
	"github.com/ambientsound/visp/input"
	"github.com/ambientsound/visp/input/keys"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/multibar"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/pkg/library"
	"github.com/ambientsound/visp/pkg/search"
	"github.com/ambientsound/visp/player"
	"github.com/ambientsound/visp/spotify/library"
	spotify_proxyclient "github.com/ambientsound/visp/spotify/proxyclient"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/ambientsound/visp/style"
	"github.com/ambientsound/visp/tabcomplete"
	"github.com/ambientsound/visp/tokencache"
	"github.com/ambientsound/visp/widgets"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify"
)

const (
	changePlayerStateDelay    = time.Millisecond * 100
	refreshInvalidTokenDeploy = time.Millisecond * 1
	refreshTokenRetryInterval = time.Second * 30
	refreshTokenTimeout       = time.Second * 5
	tickerInterval            = time.Second * 1
)

type Visp struct {
	Termui     *widgets.Application
	Tokencache tokencache.Tokencache

	client       *spotify.Client
	clipboards   *clipboard.List
	commands     chan string
	db           *db.List
	history      list.List
	index        library.Index
	interpreter  *input.Interpreter
	library      *spotify_library.List
	list         list.List
	callbacks    chan func() error
	multibar     *multibar.Multibar
	player       *player.State
	quit         chan interface{}
	sequencer    *keys.Sequencer
	stylesheet   style.Stylesheet
	ticker       *time.Ticker
	tokenRefresh <-chan time.Time
}

var _ api.API = &Visp{}

func (v *Visp) Init() {
	tcf := func(in string) multibar.TabCompleter {
		return tabcomplete.New(in, v)
	}
	v.clipboards = clipboard.New()
	v.callbacks = make(chan func() error, 16)
	v.commands = make(chan string, 1024)
	v.db = db.New()
	v.interpreter = input.NewCLI(v)
	v.library = spotify_library.New()
	v.multibar = multibar.New(tcf)
	v.player = player.NewState(spotify.PlayerState{})
	v.quit = make(chan interface{}, 1)
	v.sequencer = keys.NewSequencer()
	v.stylesheet = make(style.Stylesheet)
	v.ticker = time.NewTicker(tickerInterval)
	v.tokenRefresh = make(chan time.Time)

	v.SetList(log.List(log.InfoLevel))
}

func (v *Visp) Main(ctx context.Context) error {
	defer v.index.Close()

	// searchCancel() is called any time a search query string arrives.
	searchCtx, searchCancel := context.WithCancel(ctx)

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			log.Errorf("Killed by signal.")
			searchCancel()
			return fmt.Errorf("killed by signal")

		case <-v.quit:
			log.Infof("Exiting.")
			searchCancel()
			return nil

		case <-v.ticker.C:
			err := v.updatePlayer()
			if err != nil {
				log.Errorf("Update player: %s", err)
				if isSpotifyAccessTokenExpired(err) {
					v.tokenRefresh = time.After(refreshInvalidTokenDeploy)
				}
			}
			v.ticker.Reset(tickerInterval)

		case <-v.tokenRefresh:
			log.Infof("Spotify access token is too old, refreshing...")
			err := v.refreshToken()
			if err != nil {
				log.Errorf("Refresh Spotify access token: %s", err)
			}

		case callback := <-v.callbacks:
			err := callback()
			if err != nil {
				log.Errorf("%s", err)
			}

		// Send commands from the multibar into the main command queue.
		case command := <-v.multibar.Commands():
			v.commands <- command

		// Search input box.
		case query := <-v.multibar.Searches():
			if len(query) == 0 {
				break
			}

			searchCancel()
			searchCtx, searchCancel = context.WithCancel(ctx)

			delay := time.Millisecond * time.Duration(options.GetInt(options.SearchDelay))
			client, _ := v.Spotify()

			results := search.Delayed(query, searchCtx, delay, client)
			// listID := uuid.New().String()
			listID := "search-results"
			columns := options.GetString(options.ColumnsTracklists)

			go func() {
				i := 0
				for result := range results {
					i++
					v.callbacks <- func() error {
						result.SetID(listID)
						result.SetVisibleColumns(strings.Split(columns, ","))
						v.SetList(result)
						return nil
					}
				}
			}()

		// Process the command queue.
		case command := <-v.commands:
			err := v.Exec(command)
			if err != nil {
				log.Errorf(err.Error())
				v.multibar.Error(err)
			}

		case ev := <-v.Termui.Events():
			// First try to handle basic terminal events, such as resize.
			if v.Termui.HandleEvent(ev) {
				break
			}

			// If the input line is activated, send key events there.
			// This function handles text input, readline-like navigation, and history.
			if v.multibar.Input(ev) {
				break
			}

			// Add the key event to the sequencer, which will determine if a keybinding was pressed.
			cmd := v.keyEventCommand(ev)
			if len(cmd) == 0 {
				break
			}
			v.commands <- cmd
		}

		// Draw UI after processing any event.
		v.Termui.Draw()
	}

	searchCancel()

	return ctx.Err()
}

// Record the current "liked" status of the current track.
func (v *Visp) updateLiked() error {
	if v.player.Item == nil || len(v.player.Item.ID) == 0 {
		return nil
	}

	log.Debugf("Fetching liked status")

	client, err := v.Spotify()
	if err != nil {
		return err
	}

	liked, err := client.UserHasTracks(v.player.Item.ID)
	if err != nil {
		return err
	}

	if len(liked) != 1 {
		return nil
	}

	v.player.SetLiked(liked[0])
	log.Debugf("Likes current track: %v", v.player.Liked())

	return nil
}

func (v *Visp) updatePlayer() error {
	var err error

	now := time.Now()
	pollInterval := time.Second * time.Duration(options.GetInt(options.PollInterval))

	// no time for polling yet; just increase the ticker.
	if v.player.CreateTime.Add(pollInterval).After(now) {
		v.player.Tick()
		return nil
	}

	log.Debugf("Fetching new player information")

	client, err := v.Spotify()
	if err != nil {
		return err
	}

	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	currentID := spotify.ID(v.player.TrackRow.ID())

	v.player.Update(*state)

	// If track changed, clear information about whether this song is liked or not
	if state.Item == nil || currentID != state.Item.ID {
		v.player.ClearLiked()
	}

	// If track changed, and is known, add the currently playing track to history
	if state.Item != nil && currentID != state.Item.ID {
		v.History().Add(spotify_tracklist.FullTrackRow(*state.Item))
	}

	if v.player.LikedIsKnown() {
		return nil
	}

	err = v.updateLiked()
	if err != nil {
		return fmt.Errorf("get liked status of current song: %s", err)
	}

	return nil
}

// KeyInput receives key input signals, checks the sequencer for key bindings,
// and runs commands if key bindings are found.
func (v *Visp) keyEventCommand(event tcell.Event) string {
	ev, ok := event.(*tcell.EventKey)
	if !ok {
		return ""
	}

	contexts := commands.Contexts(v)
	v.sequencer.KeyInput(ev, contexts)
	match := v.sequencer.Match(contexts)

	if match == nil {
		return ""
	}

	log.Debugf("Input sequencer matches bind: '%s' -> '%s'", match.Sequence, match.Command)

	return match.Command
}

// SourceDefaultConfig reads, parses, and executes the default config.
func (v *Visp) SourceDefaultConfig() error {
	reader := strings.NewReader(options.Defaults)
	return v.SourceConfig(reader)
}

// SourceConfigFile reads, parses, and executes a config file.
func (v *Visp) SourceConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Infof("Reading configuration file %s", path)
	return v.SourceConfig(file)
}

// SourceConfig reads, parses, and executes config lines.
func (v *Visp) SourceConfig(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		err := v.interpreter.Exec(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Visp) refreshToken() error {
	server := options.GetString(options.SpotifyAuthServer)
	client := &http.Client{
		Timeout: refreshTokenTimeout,
	}
	token, err := spotify_proxyclient.RefreshToken(server, client, v.Tokencache.Cached())
	if err != nil {
		v.tokenRefresh = time.After(refreshTokenRetryInterval)
		return err
	}
	return v.Authenticate(token)
}

func isSpotifyAccessTokenExpired(err error) bool {
	match, _ := regexp.MatchString("access token", err.Error())
	return match
}
