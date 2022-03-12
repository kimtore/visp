package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/ambientsound/visp/list"
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

	trackAttributes = []string{
		"acousticness",
		"danceability",
		"duration",
		"energy",
		"instrumentalness",
		"key",
		"liveness",
		"loudness",
		"mode",
		"popularity",
		"speechiness",
		"tempo",
		"time_signature",
		"valence",
	}
)

// Recommend gives a list of tracks similar to the ones selected and within certain constraints.
// Effectively it implements "track radio" from the official client, but with granular control.
type Recommend struct {
	command
	api            api.API
	seedType       string
	attributes     *spotify.TrackAttributes
	usedAttributes map[string]interface{}
}

// NewRecommend returns Recommend.
func NewRecommend(api api.API) Command {
	return &Recommend{
		api:            api,
		attributes:     spotify.NewTrackAttributes(),
		usedAttributes: make(map[string]interface{}),
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
		}
	}

	if len(cmd.seedType) == 0 {
		return fmt.Errorf("wrong seed type '%s'; expected one of %s", lit, strings.Join(seedTypes, ", "))
	}

	for {
		tok, lit = cmd.Scan()
		switch tok {
		case lexer.TokenEnd:
			return nil
		case lexer.TokenWhitespace:
			break
		default:
			return fmt.Errorf("unexpected '%s'; expected track attribute or END", lit)
		}

		err := cmd.parseTrackAttribute()
		if err != nil {
			return err
		}
	}
}

func (cmd *Recommend) setTabCompleteAttributes(lit string) {
	attrs := make([]string, 0, len(trackAttributes))
	for _, attr := range trackAttributes {
		_, ok := cmd.usedAttributes[attr]
		if !ok {
			attrs = append(attrs, attr)
		}
	}
	cmd.setTabComplete(lit, attrs)
}

func (cmd *Recommend) parseTrackAttribute() error {
	tok, lit := cmd.Scan()

	cmd.setTabCompleteAttributes(lit)

	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("unexpected '%s'; expected track attribute", lit)
	}

	name := lit
	cmd.usedAttributes[lit] = new(interface{}) // remove from tab completion candidates

	tok, lit = cmd.Scan()
	if tok != lexer.TokenEqual {
		return fmt.Errorf("unexpected '%s'; expected equal sign", lit)
	}

	cmd.setTabCompleteEmpty()

	min, err := cmd.ParseUnsignedFloat()
	if err != nil {
		return err
	}

	tok, lit = cmd.Scan()
	if tok != lexer.TokenMinus {
		cmd.Unscan()
		return cmd.setTrackAttribute(cmd.attributes, name, min, nil)
	}

	max, err := cmd.ParseUnsignedFloat()
	if err != nil {
		return err
	}

	return cmd.setTrackAttribute(cmd.attributes, name, min, &max)
}

func (cmd *Recommend) setTrackAttributeInt(minFunc, maxFunc, targetFunc func(int) *spotify.TrackAttributes, min float64, max *float64) {
	if max == nil {
		targetFunc(int(min))
		return
	}
	if min > *max {
		tmp := *max
		*max = min
		min = tmp
	}
	minFunc(int(min))
	maxFunc(int(*max))
}

func (cmd *Recommend) setTrackAttributeFloat(minFunc, maxFunc, targetFunc func(float64) *spotify.TrackAttributes, min float64, max *float64) {
	if max == nil {
		targetFunc(min)
		return
	}
	if min > *max {
		tmp := *max
		*max = min
		min = tmp
	}
	minFunc(min)
	maxFunc(*max)
}

func (cmd *Recommend) setTrackAttribute(attributes *spotify.TrackAttributes, name string, min float64, max *float64) error {
	switch name {
	case "acousticness":
		cmd.setTrackAttributeFloat(attributes.MinAcousticness, attributes.MaxAcousticness, attributes.TargetAcousticness, min, max)
	case "danceability":
		cmd.setTrackAttributeFloat(attributes.MinDanceability, attributes.MaxDanceability, attributes.TargetDanceability, min, max)
	case "duration":
		// duration expected in milliseconds, translate to seconds
		min *= 1000
		if max != nil {
			*max *= 1000
		}
		cmd.setTrackAttributeInt(attributes.MinDuration, attributes.MaxDuration, attributes.TargetDuration, min, max)
	case "energy":
		cmd.setTrackAttributeFloat(attributes.MinEnergy, attributes.MaxEnergy, attributes.TargetEnergy, min, max)
	case "instrumentalness":
		cmd.setTrackAttributeFloat(attributes.MinInstrumentalness, attributes.MaxInstrumentalness, attributes.TargetInstrumentalness, min, max)
	case "key":
		cmd.setTrackAttributeInt(attributes.MinKey, attributes.MaxKey, attributes.TargetKey, min, max)
	case "liveness":
		cmd.setTrackAttributeFloat(attributes.MinLiveness, attributes.MaxLiveness, attributes.TargetLiveness, min, max)
	case "loudness":
		cmd.setTrackAttributeFloat(attributes.MinLoudness, attributes.MaxLoudness, attributes.TargetLoudness, min, max)
	case "mode":
		cmd.setTrackAttributeInt(attributes.MinMode, attributes.MaxMode, attributes.TargetMode, min, max)
	case "popularity":
		cmd.setTrackAttributeInt(attributes.MinPopularity, attributes.MaxPopularity, attributes.TargetPopularity, min, max)
	case "speechiness":
		cmd.setTrackAttributeFloat(attributes.MinSpeechiness, attributes.MaxSpeechiness, attributes.TargetSpeechiness, min, max)
	case "tempo":
		cmd.setTrackAttributeFloat(attributes.MinTempo, attributes.MaxTempo, attributes.TargetTempo, min, max)
	case "time_signature":
		cmd.setTrackAttributeInt(attributes.MinTimeSignature, attributes.MaxTimeSignature, attributes.TargetTimeSignature, min, max)
	case "valence":
		cmd.setTrackAttributeFloat(attributes.MinValence, attributes.MaxValence, attributes.TargetValence, min, max)
	default:
		return fmt.Errorf("unsupported track attribute '%s'", name)
	}

	return nil
}

func (cmd *Recommend) seeds(seedType string, tracks []list.Row) (*spotify.Seeds, error) {
	switch seedType {
	case SeedTypeArtist:
		return nil, fmt.Errorf("FIXME: seed type '%s' is unimplemented", seedType)

	case SeedTypeTrack:
		ids := make([]spotify.ID, len(tracks))
		for i := range tracks {
			ids[i] = spotify.ID(tracks[i].ID())
		}
		return &spotify.Seeds{
			Tracks: ids,
		}, nil
	default:
		return nil, fmt.Errorf("wrong seed type '%s'; expected one of %s", seedType, strings.Join(seedTypes, ", "))
	}
}

func (cmd *Recommend) name(seedType string, tracks []list.Row) string {
	var name string

	titles := make([]string, len(tracks))
	switch seedType {
	case SeedTypeArtist:
		name = "Tracks similar to those from artist(s) "
		for i := range tracks {
			titles[i] = strconv.Quote(tracks[i].Get("artist"))
		}
	case SeedTypeTrack:
		name = "Tracks similar to "
		for i := range tracks {
			titles[i] = strconv.Quote(tracks[i].Get("title"))
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
	list := cmd.api.List()
	selection := list.Selection()

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	seedTracks := selection.All()
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
	if err != nil {
		return err
	}
	fullTrackList := spotify_tracklist.NewFromTracks(fullTracks)

	err = selection.InsertList(fullTrackList, selection.Len())
	if err != nil {
		return err
	}

	selection.SetName(cmd.name(cmd.seedType, seedTracks))
	selection.SetID(uuid.New().String())
	selection.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	cmd.api.SetList(selection)
	list.ClearSelection()

	log.Infof("Copied %d source tracks and %d recommendations into '%s'", selection.Len(), len(recommendations.Tracks), selection.Name())

	return nil
}
