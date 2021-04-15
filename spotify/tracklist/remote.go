package spotify_tracklist

import (
	"time"
)

// Declare that this tracklist represents a server-side playlist,
// enabling it to be saved back to the Spotify servers.
// The playlist ID is the same as the list ID.
func (l *List) SetRemote(remote bool) {
	l.remote = remote
}

// Returns true if the tracklist represents a remote playlist.
func (l *List) HasRemote() bool {
	return l.remote
}

// Returns true if the tracklist has local changes that are not synced remotely.
func (l *List) HasLocalChanges() bool {
	return l.HasRemote() && l.Updated().After(l.syncTime)
}

// Use this function to indicate that the local and remote copies are in sync.
func (l *List) SetSyncedToRemote() {
	l.syncTime = time.Now()
}
