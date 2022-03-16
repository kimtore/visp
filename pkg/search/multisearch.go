package search

import (
	"context"
	"fmt"
	"time"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	spotify_aggregator "github.com/ambientsound/visp/spotify/aggregator"
	"github.com/zmb3/spotify/v2"
)

// Delayed enables search-as-you-type behavior.
//
// The delay parameter specifies a delay between searching the index and searching Spotify.
// If the provided context is canceled before that, the Spotify query is never performed.
// This prevents spamming the Spotify API if the user types fast.
func Delayed(query string, ctx context.Context, delay time.Duration, client *spotify.Client) <-chan list.List {
	ch := make(chan list.List, 2)

	go func() {
		defer close(ch)

		if client == nil {
			return
		}

		// 1. wait for timeout to perform spotify search, or bail out
		select {
		case <-ctx.Done():
			return
		case <-time.NewTimer(delay).C:
			break
		}

		// 2. send query to spotify
		results, err := spotify_aggregator.Search(*client, query, options.GetInt(options.Limit))
		if err != nil {
			log.Errorf("spotify search failed: %s", err)
			return
		}

		results.SetName(fmt.Sprintf("Search for '%s' (%d results)", query, results.Len()))
		ch <- results
	}()

	return ch
}
