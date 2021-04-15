package spotify_tracklist

import (
	"github.com/zmb3/spotify"
)

const maxAddToPlaylist = 100

// Add more than 100 tracks to a Spotify playlist by looping through the function.
func AddTracksToPlaylist(client *spotify.Client, playlistID spotify.ID, ids []spotify.ID) (string, error) {
	var snapshot string
	var err error

	for i := 0; i < len(ids); i += maxAddToPlaylist {
		batch := ids[i:]
		if len(batch) > maxAddToPlaylist {
			batch = batch[:maxAddToPlaylist]
		}
		snapshot, err = client.AddTracksToPlaylist(playlistID, batch...)
		if err != nil {
			break
		}
	}

	return snapshot, err
}
