package spotify_tracklist

import (
	"context"

	"github.com/ambientsound/visp/utils"
	"github.com/zmb3/spotify/v2"
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
		snapshot, err = client.AddTracksToPlaylist(context.TODO(), playlistID, batch...)
		if err != nil {
			break
		}
	}

	return snapshot, err
}

// Replace more than 100 tracks in a Spotify playlist.
func ReplacePlaylistTracks(client *spotify.Client, playlistID spotify.ID, ids []spotify.ID) error {
	batchSize := utils.Min(maxAddToPlaylist, len(ids))
	err := client.ReplacePlaylistTracks(context.TODO(), playlistID, ids[:batchSize]...)

	if err != nil || len(ids) <= maxAddToPlaylist {
		return err
	}

	_, err = AddTracksToPlaylist(client, playlistID, ids[maxAddToPlaylist:])

	return err
}
