package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
)

const (
	SeedTypeArtist = "artist"
	SeedTypeGenre  = "genre"
	SeedTypeTrack  = "track"
)

var (
	seedTypes = []string{
		SeedTypeArtist,
		// SeedTypeGenre,
		SeedTypeTrack,
	}
)

// Recommend gives a list of tracks similar to the ones selected and within certain constraints.
// Effectively it implements "track radio" from the official client, but with granular control.
type Recommend struct {
	command
	api        api.API
	seedType   string
	attributes *spotify.TrackAttributes
}

// NewRecommend returns Recommend.
func NewRecommend(api api.API) Command {
	return &Recommend{
		api:        api,
		attributes: spotify.NewTrackAttributes(),
	}
}

// Parse implements Command.
func (cmd *Recommend) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabComplete(lit, seedTypes)

	switch tok {
	case lexer.TokenIdentifier:
	case lexer.TokenEnd:
		cmd.seedType = SeedTypeTrack
		return nil
	}

	for _, seedType := range seedTypes {
		if lit == seedType {
			cmd.seedType = lit
			cmd.setTabCompleteEmpty()
			return nil
		}
	}

	return fmt.Errorf("wrong seed type '%s'; expected one of %s", lit, strings.Join(seedTypes, ", "))
}

func (cmd *Recommend) seeds(seedType string, tracks []spotify.FullTrack) (*spotify.Seeds, error) {
	switch seedType {
	case SeedTypeArtist:
		ids := make([]spotify.ID, len(tracks))
		for i := range tracks {
			ids[i] = tracks[i].Artists[0].ID
		}
		return &spotify.Seeds{
			Artists: ids,
		}, nil

	case SeedTypeTrack:
		ids := make([]spotify.ID, len(tracks))
		for i := range tracks {
			ids[i] = tracks[i].ID
		}
		return &spotify.Seeds{
			Tracks: ids,
		}, nil
	default:
		return nil, fmt.Errorf("wrong seed type '%s'; expected one of %s", seedType, strings.Join(seedTypes, ", "))
	}
}

func (cmd *Recommend) name(seedType string, tracks []spotify.FullTrack) string {
	var name string

	titles := make([]string, len(tracks))
	switch seedType {
	case SeedTypeArtist:
		name = "Tracks similar to those from artist(s) "
		for i := range tracks {
			titles[i] = strconv.Quote(tracks[i].Artists[0].Name)
		}
	case SeedTypeTrack:
		name = "Tracks similar to "
		for i := range tracks {
			titles[i] = strconv.Quote(tracks[i].Name)
		}
	}
	switch len(titles) {
	case 0:
		name += "nothing"
	case 1:
		name += titles[0]
	case 2:
		name += strings.Join(titles, " and ")
	default:
		name += strings.Join(titles[:2], ", ") + fmt.Sprintf(" and %d more", len(titles)-2)
	}

	return name
}

// Exec implements Command.
func (cmd *Recommend) Exec() error {
	list := cmd.api.Tracklist()
	if list == nil {
		return fmt.Errorf("`recommend` only works in tracklists")
	}

	selection := list.Selection()

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	seedTracks := selection.Tracks()
	list.CommitVisualSelection()
	list.DisableVisualSelection()

	seeds, err := cmd.seeds(cmd.seedType, seedTracks)
	if err != nil {
		return err
	}

	limit := options.GetInt(options.Limit)
	recommendations, err := client.GetRecommendations(*seeds, cmd.attributes, &spotify.Options{
		Limit: &limit,
	})

	if err != nil {
		return err
	}

	fullTracks, err := spotify_tracklist.SimpleTracksToFullTracks(client, recommendations.Tracks)

	newList := spotify_tracklist.NewFromTracks(fullTracks)
	newList.SetName(cmd.name(cmd.seedType, seedTracks))
	newList.SetID(uuid.New().String())
	newList.SetVisibleColumns(options.GetList(options.Columns))

	cmd.api.SetList(newList)
	list.ClearSelection()

	log.Infof("Received %d recommendations into '%s'", newList.Len(), newList.Name())

	return nil
}
