package spotify_tracklist

import (
	"github.com/zmb3/spotify/v2"
)

func artistNames(artists []spotify.SimpleArtist) []string {
	names := make([]string, len(artists))
	for i := range artists {
		names[i] = artists[i].Name
	}
	return names
}
