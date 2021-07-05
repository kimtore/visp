# Commands

Commands are strings of text, and can be entered in the [_multibar_](#switching-input-modes) or in configuration files.
Enter commands by pressing `:`. This will place you in _command input mode_. Exit command input mode by pressing `<Ctrl-C>`
or completing a command by pressing `<Enter>`.

Navigate command history by pressing `<Up>` or `<Down>`. Delete the previous word by pressing `<Ctrl-W>`.
Clear the entire line with `<Ctrl-U>`.

Commands along with their parameters can be _tab completed_ by pressing `<Tab>` at any point.
Press `<Tab>` multiple times to cycle through all available options.

Below is a list of commands recognized by Visp, along with their parameters and description.

Literal text spells out normally, placeholders enclosed in `<angle brackets>`, and optional parameters enclosed in `[square brackets]`.


## Move the cursor and viewport

These commands are split into `cursor` and `viewport` namespaces.
`cursor` commands primarily move the cursor, and `viewport` commands primarily move the viewport.

* `cursor current`

  Move the cursor to the currently playing track.

* `cursor up`  
  `cursor down`

  Move the cursor up or down one row.

* `cursor home`  
  `cursor end`

  Move the cursor to the very first or last track in the current list.

* `cursor high`  
  `cursor middle`  
  `cursor low`

  Move the cursor to the top, middle, or bottom of the current viewport.

* `cursor nextOf <tag> [<tag> [...]]`

  Move the cursor down to the next track on the list
  where any of the given tags do not match the corresponding tag of the current track.

  For example, `cursor nextOf album musicbrainz_albumid` will move to the first track of the next album;
  more specifically, the first track where either the album title or the Musicbrainz album ID does not match the corresponding tag on the current track.

* `cursor prevOf <tag> [<tag> [...]]`

  Move the cursor up to the last track on the list in sequence
  where all of the given tags match the corresponding tags of the current track.
  If the current song *was* the last match, continue searching upwards until the tags differ again.

  For example, `cursor prevOf artist musicbrainz_artistid` will move the cursor up the list
  to the first track by the current artist;
  more specifically, to the final track encountered in sequence where both of these tags match.
  If used on the top track by an artist, the cursir will move up further to the first track of the previous artist.

* `cursor random`

  Move the cursor to a random position in the current list.

* `cursor +<N>`  
  `cursor -<N>`

  Move the cursor by a particular number of rows.
  Negative numbers move the cursors up, and positive numbers move the cursor down.

* `cursor <N>`

  Move the cursor to an absolute position in the current list, where `0` is the very first track.

* `viewport up`  
  `viewport down`

  Move the viewport up or down one row;
  leave the cursor on its current song if possible.

* `viewport halfpageup`  
  `viewport halfpgup`  
  `viewport halfpaged[ow]n`  
  `viewport halfpgdn`

  Move the viewport up or down a number of rows equal to half the viewport height.
  Independently, move the cursor up or down the same number.

* `viewport pageup`  
  `viewport pgup`  
  `viewport paged[ow]n`  
  `viewport pgdn`

  Move the viewport up or down one full page (actually slightly less in most cases),
  leaving the cursor on its current song where possible.

* `viewport high`  
  `viewport low`

  Move the viewport up as high or as low as possible while leaving the cursor in view,
  still pointing to the same song.
  (When `center` is set the cursor will not end up pointing to the same song.)

* `viewport middle`

  Move the viewport so that the cursor is in the middle of the viewport,
  still pointing to the same song.


## Manipulating lists

These commands switch between, create, and edit lists.

* `list next`  
  `list prev`

  Switch to the next or previous list.

* `list <N>`

  Switch to the list with the given index.

* `list new [playlist name]`

  Create a new track list. The list remains in memory until `write` is used to save it to Spotify.
  
* `list duplicate`

  Duplicate the current list.

* `list home`

  Switch to _log console_.
  
* `list goto <id>`

  Switch to a named list. `id` can be a Spotify ID.

* `list last`

  Activate the previous window. Due to its special nature, the list of windows, nor the list of clipboards, will
  never be considered the previous window as this would result in a bad user experience.

* `list close`

  Close the currently visible list. At least one list needs to be visible, so if all lists are closed, the log console is opened.

* `isolate <tag> [<tag> [...]]`

  Search for tracks with similar tags to the current [selection](#selecting-tracks), and create a new tracklist with the results.
  The tracklist is sorted by the default sort criteria.

  See also [`inputmode search`](#switching-input-modes) for another way to create new lists.

* `sort [<tag> [...]]`

  Sort the current tracklist by the tags specified in the `sort` option if no tags are given, or otherwise by the specified tags.
  The most significant sort criterion is specified last.

  The first sort is performed as an unstable sort, while the remainder use a stable sorting algorithm.
  
* `columns <column>[ <column>[...]]`

  Specify columns that should be visible in the current list.


## Spotify library
  
* `like add cursor`  
  `like add current`  
  `like add selection`  
  `like remove cursor`  
  `like remove current`  
  `like remove selection`  
  `like toggle cursor`  
  `like toggle current`  
  `like toggle selection`
  
  Add or remove tracks from your Spotify library, also known as _liked songs_.
  Using `current` will like the track currently playing, whereas `cursor` and `selection`
  likes the track(s) under the cursor or currently selected, respectively.
  If there is no selection, `like ... selection` acts as `like ... cursor`.
  
* `recommend`  
  `recommend artist [attr=<TARGET|MIN-MAX>] [...]`  
  `recommend track [attr=<TARGET|MIN-MAX>] [...]`

  Get a list of song recommendations based on the currently selected tracks, or if no selection, the track beneath the cursor.
  `recommend artist` will make recommendations on the track artists, whereas `recommend track` picks recommendations based
  on the tracks themselves.
  
  `recommend` without any parameters behaves as `recommend track`.
  
  You may constrain the results of recommendations by specifying one or more _attributes_, along with a target value,
  or optionally floor and ceiling values of _MIN_ and _MAX_.
  
  | Attribute | Minimum value | Maximum value | Description |
  |-----------|---------------|---------------|-------------|
  | `acousticness` | 0.0 | 1.0 | A confidence measure from 0.0 to 1.0 of whether the track is acoustic. |
  | `danceability` | 0.0 | 1.0 | How suitable a track is for dancing based on a combination of musical elements including tempo, rhythm stability, beat strength, and overall regularity. |
  | `duration` | 0 | +Inf | Duration of a track, in seconds. |
  | `energy` | 0.0 | 1.0 | Perceptual measure of intensity and activity.  Typically, energetic tracks feel fast, loud, and noisy. |
  | `instrumentalness` | 0.0 | 1.0 | Instrumentalness predicts whether a track contains no vocals. "Ooh" and "aah" sounds are treated as instrumental in this context. Rap or spoken word tracks are clearly "vocal". |
  | `key` | −39 | 48 | [Pitch class notation](https://en.wikipedia.org/wiki/Pitch_class) of track root key. |
  | `liveness` | 0.0 | 1.0 | Detects the presence of an audience in the recording.  Higher liveness values represent an increased probability that the track was performed live. |
  | `loudness` | -60.0 | 0.0 | Loudness values are averaged across the entire track and are useful for comparing the relative loudness of tracks. Measured in dB. |
  | `mode` | 0.0 | 1.0 | Indicates the modality (major or minor) of a track, the type of scale from which its melodic content is derived. Major is represented by 1 and minor is 0. |
  | `popularity` | 0 | 100 | The popularity is calculated by algorithm and is based, in the most part, on the total number of plays the track has had and how recent those plays are. |
  | `speechiness` | 0.0 | 1.0 | The more exclusively speech-like the recording, the closer to 1.0 the speechiness will be. Values above 0.66 describe tracks that are probably made entirely of spoken words.  Values between 0.33 and 0.66 describe tracks that may contain both music and speech, including such cases as rap music. Values below 0.33 most likely represent music and other non-speech-like tracks.
  | `tempo` | 0 | +Inf | The overall estimated tempo of a track in beats per minute (BPM). |
  | `time_signature` | -Inf | +Inf | An estimated overall time signature of a track. The time signature (meter) is a notational convention to specify how many beats are in each bar (or measure). |
  | `valence` | 0.0 | 1.0 | Describes the musical positiveness conveyed by a track. Tracks with high valence sound more positive (e.g. happy, cheerful, euphoric), while tracks with low valence sound more negative (e.g. sad, depressed, angry). |

* `rename <playlist name>`

  Assign a new name to the current playlist. Changes must be saved back to Spotify with `w[rite]`.

* `w[rite] [playlist name]`

  Save the current track list as a Spotify playlist.
  If no name is given, and the track list is an existing Spotify playlist, Visp will save changes to this list.
  If a name is given, Visp will create a new Spotify playlist.
  

### Adding, removing, and moving tracks

* `add [<uri> [...]]`

  Add one or more files or URIs to the queue.
  If no parameters are given, the current [selection](#selecting-tracks) is assumed.

  See also [`play cursor` and `play selection`](#controlling-playback).

* `yank`  
  `copy`

  Replace the clipboard contents with the currently selected tracks.

* `cut`

  Remove the current [selection](#selecting-tracks) from the tracklist, and replace the clipboard contents with the removed tracks.

* `paste [after]`  
  `paste before`

  Insert the contents of the clipboard after (this is default) or before the cursor position.


## Selecting tracks

The `select` commands allow the tracklist selection to be manipulated.

* `select toggle`

  Toggle selection status for the track under the cursor.

  When used from visual mode, all tracks currently in the visual selection will have their manual selection status toggled, and visual mode is switched off.

* `select visual`

  Toggle _visual mode_ and anchor the selection on the track under the cursor.

* `select nearby <tag> [<tag> [...]]`

  Set the visual selection to nearby tracks with the same specified tags as the track under the cursor.
  If there is already a visual selection, it will be cleared instead.

* `select intersect <tracklist>`

  Selects all tracks present in both the current and the named list.
  The named list must be known by Visp. Known lists are either contained within the window list,
  or referred to by a list of Spotify playlists.

* `select duplicates`

  Select extra copies of all tracks found in the current list.


## Controlling playback

* `prev[ious]`

  Skip back to the previous track.

* `next`

  Skip to the next track.

* `pause`

  Pause or resume playback.

* `play`

  Start playing the queue, or resume playing the current song.

* `play cursor`

  Start playing the current list, starting at the cursor position.

* `play selection`

  Start playing the current [selection](#selecting-tracks).
  If there is no selection, fall back to `play cursor` above.

* `seek +<N>`  
  `seek -<N>`

  Seek relatively by a given number of seconds.

* `seek <N>`

  Seek to a particular point in the song, measured in seconds.

* `stop`

  Stop playback.

* `repeat`  
  `repeat context`  
  `repeat track`
  `repeat off`

  Switch between repeat modes. `context` will repeat the currently playing list or selection, whereas `track` will
  repeat the currently playing track. Running this command without parameters will toggle between the different modes.

* `shuffle`  
  `shuffle on`
  `shuffle off`

  Switch between shuffle modes. Running this command without parameters will toggle shuffle on and off.

### Controlling the volume

These commands control the volume. The volume range is from 0 to 100.

* `volume <N>`

  Set the volume to an absolute value.

* `volume +<N>`  
  `volume -<N>`

  Adjust the volume by a relative value.

* `volume mute`

  Toggle mute.


## Switching input modes

* `inputmode normal`

  Switch to normal mode, where key bindings take effect.

* `inputmode input`

  Switch to input mode: focus the multibar, where commands can be typed in.

* `inputmode search`

  Switch to search mode, where searches execute as you type.

  When `<Enter>` is pressed from search mode, the result is a new list containing the current search results.


## Customizing Visp

### Setting global options

The command `set` or its shorthand `se` can be used to change global program options at runtime.
The [list of available options](options.md) is documented elsewhere.

* `set <option>=<value>`

  Set a non-boolean option to a particular value.

* `set <option>`  
  `set no<option>`

  Switch a boolean option on or off.

* `set inv<option>`  
  `set <option>!`

  Toggle a boolean option.

* `set <option>?`

  Query the current value of an option.

### Setting key sequences

These commands bind and unbind key sequences to commands.

A _key sequence_ can have any number of elements, each of which is any of:

* a letter
* a "special" key enclosed in angle brackets, such as `<space>` or `<f1>`
* a key with modifiers, enclosed in angle brackets, such as `<Ctrl-X>`, `<Alt-A>`, or `<Shift-Escape>`

Modifier keys are `Ctrl`, `Alt`, `Meta`, and `Shift`.
The first letter of these four words can also be used as a shorthand.

Special keys are too numerous to list here, but can be found in the [complete list of special keys](/keysequence/names.go).

Regular keys such as letters, numbers, symbols, unicode characters, etc. will never have the `Shift` modifier key.
Generally, terminal applications have far less insight into keyboard activity than graphical applications,
and therefore you should avoid depending too much on availability of modifiers or any specific keys.

_Contexts_ are a way to make keybindings context sensitive. Choose between `global`, `library`, `tracklist`, `devices`, and `windows`.
You can bind a key sequence to multiple contexts. The local context takes precedence, so a sequence bound to
the `tracklist` context will always be attempted before `global`.

* `bind <context> <key sequence> <command>`

  Configure a specific keyboard input sequence to execute a command.

* `unbind <context> <key sequence>`

  Unbind a key sequence.

### Setting styles

* `style <name> [<foreground> [<background>]] [bold] [underline] [reverse] [blink]`

  Specify the style of a UI item.
  See the [styling guide](styling.md#text-style) for details.

  The keywords `bold`, `underline`, `reverse`, and `blink` can be specified literally.
  Any keyword order is accepted, but the background color, if specified, must come after the foreground color.


## Miscellaneous

* `print [<tag> [...]]`

  Show the contents of the given tag(s) for the track under the cursor.

* `q[uit]`

  Exit the program. Any unsaved changes will be lost.

* `redraw`

  Force a screen redraw. Useful if rendering has gone wrong.

* `show history`
  `show logs`  
  `show library`
  `show keybindings`
  `show windows`

  Switch between different views.
