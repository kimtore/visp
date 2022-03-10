// package spotify_library provides access to pre-defined Spotify content
// such as playlists, library, top artists, categories, and so on.
//
// This package also provides access to internal functions such as devices and clipboard.
package spotify_library

import (
	"github.com/ambientsound/visp/list"
)

type List struct {
	list.Base
}

var _ list.List = &List{}

const (
	listName = "description"
)

const (
	Categories          = "categories"
	Devices             = "devices"
	FeaturedPlaylists   = "featured-playlists"
	FollowedArtists     = "followed-artists"
	MyAlbums            = "my-albums"
	MyFollowedPlaylists = "my-followed-playlists"
	MyPlaylists         = "my-playlists"
	MyTracks            = "my-tracks"
	NewReleases         = "new-releases"
	TopArtists          = "top-artists"
	TopTracks           = "top-tracks"
)

var entries = map[string]string{
	Devices:           "Player devices",
	FeaturedPlaylists: "Featured playlists",
	NewReleases:       "New releases",
	MyPlaylists:       "Playlists from my Spotify library",
	MyTracks:          "All liked songs from my library",
	MyAlbums:          "All saved albums in my library",
	TopTracks:         "Top tracks from my listening history",
}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("spotify_library")
	this.SetName("Libraries and discovery")
	this.SetVisibleColumns([]string{listName})
	for key, name := range entries {
		this.Add(list.NewRow(
			key,
			list.DataTypeFIXME,
			map[string]string{
				listName: name,
			}))
	}
	this.Sort([]string{listName})
	this.SetCursor(0)
	return this
}
