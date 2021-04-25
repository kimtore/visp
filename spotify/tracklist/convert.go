package spotify_tracklist

import (
	"fmt"

	"github.com/zmb3/spotify"
)

const maxFullTracks = 50

// Convert a list of SimpleTrack objects to FullTrack.
func SimpleTracksToFullTracks(client *spotify.Client, simpleTracks []spotify.SimpleTrack) ([]spotify.FullTrack, error) {
	ids := make([]spotify.ID, len(simpleTracks))
	for i := range simpleTracks {
		ids[i] = simpleTracks[i].ID
	}

	var err error
	allTracks := make([]spotify.FullTrack, 0, len(simpleTracks))

	for i := 0; i < len(ids); i += maxFullTracks {
		batch := ids[i:]
		if len(batch) > maxFullTracks {
			batch = batch[:maxFullTracks]
		}
		tracks, err := client.GetTracks(batch...)
		if err != nil {
			return nil, err
		}

		for i, track := range tracks {
			if track == nil {
				return nil, fmt.Errorf("get full track information for '%s' failed", batch[i])
			}
			allTracks = append(allTracks, *track)
		}
	}

	return allTracks, err
}
