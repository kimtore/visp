package topbar

import (
	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/version"
)

// Version draws the short name of this application, as defined in the version module.
type Version struct {
	version string
}

// NewVersion returns Version.
func NewVersion(a api.API, param string) Fragment {
	return &Version{version.Version}
}

// Text implements Fragment.
func (w *Version) Text() (string, string) {
	return w.version, `version`
}
