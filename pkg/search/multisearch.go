package search

import (
	"context"
	"fmt"
	"time"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/options"
	"github.com/ambientsound/visp/pkg/library"
	spotify_aggregator "github.com/ambientsound/visp/spotify/aggregator"
	"github.com/zmb3/spotify"
)

// Multisearch enables search-as-you-type behavior by presenting a hybrid between the search index and Spotify APIs.
//
// The delay parameter specifies a delay between searching the index and searching Spotify.
// If the provided context is canceled before that, the Spotify query is never performed.
// This prevents spamming the Spotify API if the user types fast.
func Multisearch(query string, ctx context.Context, delay time.Duration, client *spotify.Client, index library.Index) <-chan list.List {
	ch := make(chan list.List, 2)

	go func() {
		defer close(ch)

		// 1. query the index first
		lst, err := index.Query(query)
		if err != nil {
			log.Errorf("index query failed: %s", err)
		} else {
			// lst.SetName(fmt.Sprintf("Search for '%s'...", query))
			// ch <- lst
		}

		// spotify client is strictly not needed, we can be content with the index query
		if client == nil {
			return
		}

		// 2. wait for timeout to perform spotify search, or bail out
		select {
		case <-ctx.Done():
			return
		case <-time.NewTimer(delay).C:
			break
		}

		// 3. send an async query to spotify
		tracklist, err := spotify_aggregator.Search(*client, query, options.GetInt(options.Limit))
		if err != nil {
			log.Errorf("spotify query failed: %s", err)
			return
		}

		if ctx.Err() != nil {
			return
		}

		// 4. append Spotify's search results
		err = lst.InsertList(tracklist, lst.Len())
		if err != nil {
			panic(err)
		}

		lst.SetName(fmt.Sprintf("Search for '%s' (%d results)", query, lst.Len()))
		ch <- lst
	}()

	return ch
}
