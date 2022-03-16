package spotify_tracklist

import (
	"context"
	"fmt"
	"strings"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/options"
	spotify_albums "github.com/ambientsound/visp/spotify/albums"
	"github.com/ambientsound/visp/utils"
	"github.com/zmb3/spotify/v2"
)

type List struct {
	list.Base
}

var _ list.List = &List{}

func NewFromFullTrackPage(client spotify.Client, source *spotify.FullTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		tracks = append(tracks, source.Tracks...)
		err = client.NextPage(context.TODO(), source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromSimpleTrackPageAndAlbum(client spotify.Client, source *spotify.SimpleTrackPage, album spotify.SimpleAlbum) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, AlbumTrack(track, album))
		}
		err = client.NextPage(context.TODO(), source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromSavedTrackPage(client spotify.Client, source *spotify.SavedTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, track.FullTrack)
		}
		err = client.NextPage(context.TODO(), source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromPlaylistTrackPage(client spotify.Client, source *spotify.PlaylistTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, track.Track)
		}
		err = client.NextPage(context.TODO(), source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil

}

func AlbumTrack(track spotify.SimpleTrack, album spotify.SimpleAlbum) spotify.FullTrack {
	return spotify.FullTrack{
		SimpleTrack: track,
		Album:       album,
	}
}

func NewFromSimpleAlbumPage(client spotify.Client, source *spotify.SimpleAlbumPage) (*List, error) {
	var lst *List
	var err error
	var trackPage *spotify.SimpleTrackPage

	tracks := make([]spotify.FullTrack, 0, source.Total)
	albums, err := spotify_albums.NewFromSimpleAlbumPage(client, source)
	if err != nil {
		return nil, err
	}

	for i := 0; i < albums.Len(); i++ {
		album := albums.Album(i)
		trackPage, err = client.GetAlbumTracks(context.TODO(), album.ID)
		if err != nil {
			break
		}
		lst, err = NewFromSimpleTrackPageAndAlbum(client, trackPage, *album)
		if err != nil {
			break
		}
		tracks = append(tracks, lst.Tracks()...)
	}

	if err != nil && err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromTracks(tracks []spotify.FullTrack) *List {
	this := &List{}
	this.Clear()
	for _, track := range tracks {
		this.Add(FullTrackRow(track))
	}
	return this
}

func NewHistory() *List {
	this := &List{}
	this.Clear()
	this.SetName("History")
	this.SetID("history")
	this.SetVisibleColumns(options.GetList(options.ColumnsTracklists))
	return this
}

func FullTrackRow(track spotify.FullTrack) list.Row {
	return &Row{
		track: track,
		BaseRow: list.NewRow(
			track.ID.String(),
			list.DataTypeTrack,
			map[string]string{
				"album":       track.Album.Name,
				"albumArtist": strings.Join(artistNames(track.Album.Artists), ", "),
				"artist":      strings.Join(artistNames(track.Artists), ", "),
				"date":        track.Album.ReleaseDateTime().Format("2006-01-02"),
				"time":        utils.TimeString(track.Duration / 1000),
				"title":       track.Name,
				"track":       fmt.Sprintf("%02d", track.TrackNumber),
				"disc":        fmt.Sprintf("%d", track.DiscNumber),
				"popularity":  fmt.Sprintf("%1.2f", float64(track.Popularity)/100),
				"year":        track.Album.ReleaseDateTime().Format("2006"),
			},
		),
	}
}

func (l *List) CursorTrack() *spotify.FullTrack {
	row := l.CursorRow()
	if row == nil {
		return nil
	}
	track := row.(*Row).track
	return &track
}

// SelectionAsTracklist returns all the selected songs as a new track list.
func (l *List) SelectionAsTracklist() *List {
	indices := l.SelectionIndices()
	tracks := make([]spotify.FullTrack, len(indices))

	for i, index := range indices {
		tracks[i] = l.Row(index).(*Row).Track()
	}

	return NewFromTracks(tracks)
}

func (l *List) Tracks() []spotify.FullTrack {
	tracks := make([]spotify.FullTrack, l.Len())
	for i := 0; i < l.Len(); i++ {
		tracks[i] = l.Row(i).(*Row).Track()
	}
	return tracks
}
