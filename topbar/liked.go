package topbar

import (
	"github.com/ambientsound/visp/api"
)

var likeString = map[bool]string{
	true:  "<3",
	false: "",
}

var likeStringUnicode = map[bool]string{
	true:  "\u2661",
	false: "",
}

// State draws the current player state as an ASCII symbol.
type Liked struct {
	api    api.API
	string map[bool]string
}

// NewState returns State.
func NewLiked(a api.API, param string) Fragment {
	str := likeString
	if param == "unicode" {
		str = likeStringUnicode
	}
	return &Liked{a, str}
}

// Text implements Fragment.
func (w *Liked) Text() (string, string) {
	return w.string[w.api.PlayerStatus().Liked()], `liked`
}
