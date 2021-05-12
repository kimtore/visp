# User guide

Welcome to _Visp_, a Vi-like Spotify client for the terminal.

First off, Visp is not a _player_ and does not output any sound. To actually play back
music, you need a dedicated program such as either the official Spotify client, or some
other playback-capable software such as [spotifyd](https://github.com/Spotifyd/spotifyd).
This program needs to be running and connected to Spotify so that Visp can recognize it.

To authorize Visp to your Spotify account, please visit https://visp.site/authorize
to get an access token. Enter your access token with `:auth <TOKEN>`.

Visp is based around the concept of _lists_. Every view in Visp is a list of some kind.
Upon starting the program, you are shown the _log console_, which keeps track of things
that happen within the program. Other lists contain playlists, albums, or tracks.

To get an overview of all the lists you've visited while running Visp, press `w`.

To enter a command, type `:` followed by the command, then press `<Enter>`. While entering
a command, you can press `<Tab>` to engage tab completion, which will complete the word
you're typing, provided that Visp can guess what you're trying to do. Tab completion in
Visp is very powerful and can recognize pretty much anything you want to accomplish.
See [command documentation](commands.md) for a list of all supported commands.

Press `<F1>`, or enter the command `show keybindings`, at any time to show a
list of key bindings.


## Basic movement

The default bindings for movement are similar to Vim.

`j` and `k` (or `<Down>` and `<Up>` and move down and up,
`gt` and `gT` (or just `t` and `T`) move forward and back between lists.
Use `<Ctrl-F>` and `<Ctrl-B>` to move a page down or up,
or `<Ctrl-D>` and `<Ctrl-U`> for half a page at a time.
`gg` and `G` go to the very top and bottom of the list,
while `H`, `M`, and `L` go to the top, middle, and bottom of the current viewport.


## Find and play music

Your Spotify library can be accessed by pressing `c`.
Navigate to any entry and press `<Enter>` to load that list from Spotify.
A list from Spotify contain playlists, albums, tracks, or special items
such as playback devices.

Type `/` to start a new search. Type in your search query, followed by `<Enter>`.
A new tracklist will appear with search results.

Once in a _playlist view_, select a playlist and press `<Enter>` to load it into Visp.
A new window will be created in _tracklist view_, showing all the tracks in that playlist.

When you are in a tracklist view, press `<Enter>` on any track to start playing
from that track. Press `a` to add the track to the queue instead of playing right away.

Press `o` to like or unlike the currently playing track.

Select up to five tracks and type `R` to have Spotify recommend similar tracks.

To see all the tracks you've played while running Visp, type `<Ctrl-W>h`.


## Select multiple tracks

Press `m` to select a single track. The track will be highlighted in blue. To unselect the track,
move the cursor over it and press `m` again.

Pressing `v` or `V` starts _visual mode_, which starts a selection from the cursor position.
Move up or down to extend the selection. Press `v` or `V` again to exit visual mode. Any tracks
selected through visual mode will be unselected. To permanently select the tracks instead, press
`m` instead of exiting through `v` or `V`.

`<Ctrl-A>` selects all tracks, while `<Ctrl-C>` unselects all tracks.


## Cut, copy, paste

Press `x` to cut, `y` to copy, and `p` to paste. `P` pastes
_before_ the cursor, while `p` pastes _after_ the cursor.

Anything you copy or cut will be placed in a new _clipboard_. Press `C` to view all clipboards
created in your Visp session. A clipboard acts like a track list.


## Make or change a playlist

Press `<Ctrl-W>c` to create a new playlist.
This playlist exists entirely in Visp and is not saved to Spotify yet.

To rename a playlist, enter the command `rename` followed by a playlist name.

Add tracks by copying and pasting them from either a search result or another tracklist.

To save a playlist to Spotify, enter the command `write` (or `w`) optionally followed by a playlist name.
If entered without a name, `write` will save changes to an existing playlist. If you provide a name,
a new playlist is created instead. The `write` command can be used from any track list or search result.


## Advanced navigation

Sort any list by pressing `<Ctrl-S>` (default sort order), or optionally by entering the command `sort`
followed by a space-delimited list of sort keys.
Sort terms are ordered by most significant field last.

Press `<Ctrl-J>` (or type `:isolate artist`) to search for tracks with the same artist, 
or `<Ctrl-T>` (`:isolate albumartist album`) to search for tracks in the same album.

Press `&` (`:select nearby albumartist album`) to select an entire album.

Type `b` or `e` to move to the previous or next album in a tracklist.


## Configuration

By default, Visp tries to read your configuration from
`$HOME/.config/visp/visp.conf`.
If you defined paths in either `$XDG_CONFIG_DIRS` or `$XDG_CONFIG_HOME`, Visp will look for `visp.conf` there.

```
# Sample Visp configuration file.
# All whitespace, newlines, and text after a hash sign will be ignored.

# The 'center' option will make sure the cursor is always centered on screen.
set center

# Some custom keyboard bindings.
bind <Alt-Left> cursor prevOf year    # jump to previous year.
bind <Alt-Right> cursor nextOf year   # jump to next year.

# Pink statusbar.
style statusbar black darkmagenta

# Minimalistic topbar.
set topbar="Now playing: ${tag|artist} \\- ${tag|title} (${elapsed})"
```

See also [default configuration](../options/options.go).
