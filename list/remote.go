package list

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

type Remote interface {
	SetRemote(remote bool)
	HasRemote() bool
	HasLocalChanges() bool
	SetSyncedToRemote()
	URI() *spotify.URI
	SetURI(uri spotify.URI)
}

// Declare that this tracklist represents a server-side playlist,
// enabling it to be saved back to the Spotify servers.
// The playlist ID is the same as the list ID.
func (s *Base) SetRemote(remote bool) {
	s.remote = remote
}

// Returns true if the tracklist represents a remote playlist.
func (s *Base) HasRemote() bool {
	return s.remote
}

// Returns true if the tracklist has local changes that are not synced remotely.
func (s *Base) HasLocalChanges() bool {
	return s.HasRemote() && s.Updated().After(s.syncTime)
}

// Use this function to indicate that the local and remote copies are in sync.
func (s *Base) SetSyncedToRemote() {
	s.syncTime = time.Now()
}

func (s *Base) URI() *spotify.URI {
	if len(s.uri) > 0 {
		uri := s.uri
		return &uri
	}
	return nil
}

func (s *Base) SetURI(uri spotify.URI) {
	s.uri = uri
}
