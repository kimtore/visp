package topbar

import (
	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/version"
)

// Shortname draws the short name of this application, as defined in the version module.
type Shortname struct {
	shortname string
}

// NewShortname returns Shortname.
func NewShortname(a api.API, param string) Fragment {
	return &Shortname{version.ShortName}
}

// Text implements Fragment.
func (w *Shortname) Text() (string, string) {
	return w.shortname, `shortName`
}
