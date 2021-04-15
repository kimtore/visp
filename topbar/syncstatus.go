package topbar

import (
	"github.com/ambientsound/visp/api"
)

// SyncStatus shows a symbol if the playlist is not synced remotely.
type SyncStatus struct {
	api api.API
}

// NewMode returns Mode.
func NewSyncStatus(a api.API, param string) Fragment {
	return &SyncStatus{a}
}

// Text implements Fragment.
func (w *SyncStatus) Text() (string, string) {
	tracklist := w.api.Tracklist()
	if tracklist == nil || !tracklist.HasLocalChanges() {
		return "", ""
	}
	return "[+]", `syncStatus`
}
