package options

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Option names.
const (
	Center            = "center"
	Columns           = "columns"
	ColumnsPlaylists  = "columns.playlists"
	Device            = "device"
	ExpandColumns     = "expandcolumns"
	FullHeaderColumns = "fullheadercolumns"
	Limit             = "limit"
	LogFile           = "logfile"
	LogOverwrite      = "logoverwrite"
	PollInterval      = "pollinterval"
	Sort              = "sort"
	SpotifyAuthServer = "spotifyauthserver"
	Topbar            = "topbar"
)

// Option types.
const (
	boolType   = false
	intType    = 0
	stringType = ""
)

var (
	v = viper.NewWithOptions(viper.KeyDelimiter("::"))
)

// Initialize option types.
// Default values must be defined in the Defaults string.
func init() {
	v.Set(Center, boolType)
	v.Set(Columns, stringType)
	v.Set(ColumnsPlaylists, stringType)
	v.Set(Device, stringType)
	v.Set(ExpandColumns, stringType)
	v.Set(FullHeaderColumns, stringType)
	v.Set(Limit, intType)
	v.Set(LogFile, stringType)
	v.Set(LogOverwrite, boolType)
	v.Set(PollInterval, intType)
	v.Set(Sort, stringType)
	v.Set(SpotifyAuthServer, stringType)
	v.Set(Topbar, stringType)
}

// Methods for getting options from Viper.
var (
	Get       = v.Get
	GetString = v.GetString
	GetInt    = v.GetInt
	GetBool   = v.GetBool
	Set       = v.Set
	AllKeys   = v.AllKeys
)

// Split a string option into a comma-delimited list.
func GetList(key string) []string {
	return strings.Split(v.GetString(key), ",")
}

// Return a human-readable representation of an option.
// This string can be used in a config file.
func Print(key string, opt interface{}) string {
	switch v := opt.(type) {
	case string:
		return fmt.Sprintf("%s=\"%s\"", key, v)
	case int:
		return fmt.Sprintf("%s=%d", key, v)
	case bool:
		if !v {
			return fmt.Sprintf("no%s", key)
		}
		return fmt.Sprintf("%s", key)
	default:
		return fmt.Sprintf("%s=%v", key, v)
	}
}

// Default configuration file.
const Defaults string = `
# Global options
set columns.playlists=name,tracks,owner,public,collaborative
set columns=artist,title,track,album,year,time,popularity
set expandcolumns=logMessage,description,deviceName,name,artist,title,album
set fullheadercolumns=logLevel,public,collaborative,deviceName,track,year,time
set limit=50
set nocenter
set pollinterval=10
set sort=track,disc,album,year,albumArtist
set spotifyauthserver="https://visp.site"
set topbar="${tag|artist} - ${tag|title} $liked|$shortname $version|$elapsed $state $time;\\#${tag|track} ${tag|album}|${list|title} [${list|index}/${list|total}] ${synced}|$device $mode $volume;;"

# Logging
set nologoverwrite
set logfile=

# Song tag styles
style album teal
style albumArtist teal
style artist yellow
style date default
style year default
style disc default
style popularity dim
style time green
style title white
style track default
style _id gray

# Tracklist styles
style currentSong black yellow
style cursor black white
style header gray dim bold
style selection gray blue

# Key binding styles
style context teal
style keyBinding white
style command yellow

# Playlist library styles
style public green
style collaborative green
style owner teal
style name white
style tracks default

# Library styles
style description white

# Topbar styles
style deviceName teal
style deviceType teal
style elapsedTime teal
style elapsedPercentage teal
style listIndex teal
style listTitle white
style listTotal teal
style mute red
style shortName yellow
style state default
style switches teal
style syncStatus red dim
style tagMissing red
style topbar darkgray
style version gray dim
style volume green
style liked green

# Other styles
style commandText default
style currentDevice white green
style errorText black red
style logLevel dim gray
style logMessage dim gray
style readout default
style searchText white bold
style sequenceText teal
style statusbar default
style timestamp teal
style visualText teal

# Keyboard bindings: cursor and viewport movement
bind global <Up> cursor up
bind global k cursor up
bind global <Down> cursor down
bind global j cursor down
bind global <PgUp> viewport pgup
bind global <PgDn> viewport pgdn
bind global <C-b> viewport pgup
bind global <C-f> viewport pgdn
bind global <C-u> viewport halfpgup
bind global <C-d> viewport halfpgdn
bind global <C-y> viewport up
bind global <C-e> viewport down
bind global <Home> cursor home
bind global gg cursor home
bind global <End> cursor end
bind global G cursor end
bind global gc cursor current
bind global H cursor high
bind global M cursor middle
bind global L cursor low
bind global zb viewport high
bind global z- viewport high
bind global zz viewport middle
bind global z. viewport middle
bind global zt viewport low
bind global z<Enter> viewport low

# Tracklist specifics
bind tracklist b cursor prevOf album
bind tracklist e cursor nextOf album
bind tracklist R recommend track

# Keyboard bindings: input mode
bind global : inputmode input
bind global / inputmode search
bind global <F3> inputmode search
bind global v select visual
bind global V select visual

# Keyboard bindings: player and mixer
bind tracklist <Enter> play selection
bind tracklist a add
bind global <Space> pause
bind global s stop
bind global h previous
bind global l next
bind global + volume +2
bind global - volume -2
bind global <left> seek -5
bind global <right> seek +5
bind global <Alt-M> volume mute
bind global S single

# Special windows
bind global c show library
bind global C show clipboards
bind global w show windows
bind global <F1> show keybindings
bind windows <Enter> show selected
bind library <Enter> show selected
bind clipboards <Enter> show selected
bind devices <Enter> device activate
bind playlists <Enter> list open

# Keyboard bindings: other
bind global <C-l> redraw
bind global <C-s> sort
bind tracklist i print file
bind global gt list next
bind global gT list previous
bind global t list next
bind global T list previous
bind global <C-w>d list duplicate
bind global <C-g> list remove
bind global <Tab> list last
bind tracklist <C-j> isolate artist
bind tracklist <C-t> isolate albumArtist album
bind tracklist & select nearby albumArtist album
bind global m select toggle
bind global <C-a> select all
bind global <C-c> select none
bind tracklist <Delete> cut
bind tracklist x cut
bind tracklist y yank
bind global Y yank current
bind tracklist p paste after
bind tracklist P paste before
bind global o like toggle current
`
