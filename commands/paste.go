package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	spotify_tracklist "github.com/ambientsound/visp/spotify/tracklist"
)

// Paste inserts songs from the clipboard.
type Paste struct {
	command
	api      api.API
	position int
	list     *spotify_tracklist.List
}

// NewPaste returns Paste.
func NewPaste(api api.API) Command {
	return &Paste{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Paste) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.list = cmd.api.Tracklist()
	if cmd.list == nil {
		return fmt.Errorf("`paste` only works in tracklists")
	}

	cmd.setTabCompleteVerbs(lit)

	// Expect either "before" or "after".
	switch tok {
	case lexer.TokenIdentifier:
		switch lit {
		case "before":
			cmd.position = 0
		case "after":
			cmd.position = 1
		default:
			return fmt.Errorf("unexpected '%s', expected position", lit)
		}
		cmd.setTabCompleteEmpty()
		return cmd.ParseEnd()

	// Fall back to "after" if no arguments given.
	case lexer.TokenEnd:
		cmd.position = 1

	default:
		return fmt.Errorf("unexpected '%s', expected position", lit)
	}

	return nil
}

// Exec implements Command.
func (cmd *Paste) Exec() error {
	cursor := cmd.list.Cursor()
	clipboard := cmd.api.Clipboards().Active()

	if clipboard == nil {
		return fmt.Errorf("no clipboard, try `cut` or `yank` first")
	}

	ln := clipboard.Len()
	err := cmd.list.InsertList(clipboard, cursor+cmd.position)

	if err != nil {
		return err
	}

	// move cursor to position of inserted items
	// if items were inserted _before_ the cursor, it is already at the correct spot
	if cmd.position == 1 {
		cmd.list.MoveCursor(1)
	}

	cmd.api.ListChanged()

	cmd.api.Message("%d tracks inserted", ln)

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Paste) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"after",
		"before",
	})
}
