# Options

The default configuration can be seen in the [options.go source file](../options/options.go).

See the documentation on [setting options](commands.md#setting-global-options) for more information on syntax.

## Spotify

### Default playback device

* `set device=<name>`

  When starting playback, and in case no Spotify playback device is currently active,
  Visp will try to play music on the device with the given name.

### Search results limit

* `set limit=50`

  Limit the number of search results returned from Spotify.

  Lowering this number might decrease latency and will lower bandwidth usage.

### Polling interval

* `set pollinterval=10`

  The Spotify Web API offers no way to get automatically notified when the player status changes.
  Thus, polling is neccessary. The default setting will poll Spotify every ten seconds to check for
  player updates.

  When a song finishes playing, or a command against Spotify is performed,
  a poll will be made regardless of this setting.

### Authentication

* `set spotifyauthserver=https://visp.site`  

  Required in order to authenticate with Spotify. Override this setting if
  setting up your own authentication proxy server, as detailed in the
  [Spotify section](spotify.md).


## Logging

### Log file

* `set logoverwrite`  
  `set nologoverwrite`

  If set, the log file is truncated when opened. Defaults to false.
  To have any effect, this option must be set before `logfile`.

* `set logfile=/path/to/debug.log`

  Writes debugging information to a file. Logging is disabled by default.
  Setting this option or changinig the file name will write the entire log to that file.
  Be careful to set `logoverwrite` or `nologoverwrite` before enabling this option.


## Visual options

### Cursor position

* `set center`  
  `set nocenter`

  If set, the viewport is automatically moved so that the cursor stays in the center, if possible.

### Visible columns

* `set columns=<tag>[,<tag>[...]]`

  Define which tags should be shown in the tracklist.

  A comma-separated list of tag names must be given, such as the default `artist,track,title,album,year,time,popularity`.

* `set columns.playlists=<tag>[,<tag>[...]]`

  Define which tags should be shown when showing a list of playlists.
  
### Sort order

* `set sort=<tag>[,<tag>[...]]`

  Set the default sort order, for when using the [`sort` command](commands.md#manipulating-lists) without any parameters.

  A comma-separated list of tag names must be given, such as the default `track,disc,album,year,albumArtist`.

### Information bar ("top bar")

* `set topbar=<spec>`

  Define the layout and visible items in the _top bar_.
  See the [styling guide](styling.md#top-bar) for information on how to configure the top bar.

  The default value is `"${tag|artist} - ${tag|title}|$shortname $version|$elapsed $state $time;\\#${tag|track} ${tag|album}|${list|title} [${list|index}/${list|total}]|$device $mode $volume;;"`
