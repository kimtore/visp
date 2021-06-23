package commands

import (
	"fmt"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/spotify/aggregator"
	"github.com/zmb3/spotify"
)

// Select manipulates song selection within a songlist.
type Select struct {
	command
	api           api.API
	all           bool
	intersect     bool
	intersectList list.List
	none          bool
	toggle        bool
	visual        bool
	duplicates    bool
	nearby        []string
}

// NewSelect returns Select.
func NewSelect(api api.API) Command {
	return &Select{
		api:    api,
		nearby: make([]string, 0),
	}
}

// Parse implements Command.
func (cmd *Select) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case "all":
		cmd.all = true
	case "none":
		cmd.none = true
	case "toggle":
		cmd.toggle = true
	case "visual":
		cmd.visual = true
	case "duplicates":
		cmd.duplicates = true
	case "intersect":
		return cmd.parseIntersect()
	case "nearby":
		return cmd.parseNearby()
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Select) Exec() error {
	lst := cmd.api.List()

	switch {
	case cmd.intersect:
		targets := cmd.intersectList.All()
		found := 0

		for _, target := range targets {
			rownum, err := lst.RowNum(target.ID())
			if err == nil {
				lst.SetSelected(rownum, true)
				found++
			}
		}

		log.Infof("Selected %d tracks from '%s'", found, cmd.intersectList.Name())

	case cmd.toggle && lst.HasVisualSelection():
		lst.CommitVisualSelection()
		lst.DisableVisualSelection()

	case cmd.visual:
		lst.ToggleVisualSelection()
		return nil

	case len(cmd.nearby) > 0:
		return cmd.selectNearby()

	case cmd.duplicates:
		lst.ClearSelection()
		seen := make(map[string]bool)
		dupes := 0
		for i, row := range lst.All() {
			if seen[row.ID()] {
				lst.SetSelected(i, true)
				dupes++
			}
			seen[row.ID()] = true
		}
		log.Infof("Selected %d duplicates", dupes)
		return nil

	case cmd.all:
		lst.DisableVisualSelection()
		for i := 0; i < lst.Len(); i++ {
			lst.SetSelected(i, true)
		}
		return nil

	case cmd.none:
		lst.ClearSelection()
		return nil

	default:
		index := lst.Cursor()
		selected := lst.Selected(index)
		lst.SetSelected(index, !selected)
	}

	lst.MoveCursor(1)

	return nil
}

func (cmd *Select) parseIntersect() error {
	lit := cmd.ScanRemainderAsIdentifier()
	cmd.setTabComplete(lit, cmd.api.Db().Names())

	for {
		result := cmd.api.Db().Lookup(lit)
		switch typed := result.(type) {
		case list.List:
			cmd.intersect = true
			cmd.intersectList = typed
			cmd.setTabCompleteEmpty()
			return cmd.ParseEnd()
		case spotify.ID:
			err := cmd.load(typed)
			if err != nil {
				return fmt.Errorf("load '%s' from Spotify: %w", typed, err)
			}
			continue
		case nil:
			return fmt.Errorf("no such list: '%s'", lit)
		default:
			return fmt.Errorf("BUG: unknown list type '%T': %v", result, typed)
		}
	}
}

// parseNearby parses tags and inserts them in the nearby list.
func (cmd *Select) parseNearby() error {

	// Data initialization and sanity checks
	list := cmd.api.List()
	row := list.CursorRow()
	if row == nil {
		return nil
	}

	// Retrieve a list of songs
	tags, err := cmd.ParseTags(row.Keys())
	if err != nil {
		return err
	}

	cmd.nearby = tags
	return nil
}

// selectNearby selects tracks near the cursor with similar tags.
func (cmd *Select) selectNearby() error {
	list := cmd.api.List()
	index := list.Cursor()
	row := list.CursorRow()

	// In case the list has a visual selection, disable that selection instead.
	if list.HasVisualSelection() {
		list.DisableVisualSelection()
		return nil
	}

	if row == nil {
		return fmt.Errorf("can't select nearby rows; list is empty")
	}

	// Find the start and end positions
	start := list.NextOf(cmd.nearby, index+1, -1)
	end := list.NextOf(cmd.nearby, index, 1) - 1

	// Set visual selection and move cursor to end of selection
	list.SetVisualSelection(start, end, start)
	list.SetCursor(end)

	return nil
}

func (cmd *Select) load(id spotify.ID) error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	// fixme: global const
	const limit = 50

	lst, err := spotify_aggregator.ListWithID(*client, id.String(), limit)
	if err != nil {
		return err
	}

	log.Debugf("auto-loaded list '%s' from Spotify", lst.Name())

	cmd.api.Db().Cache(lst)

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Select) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"all",
		"duplicates",
		"intersect",
		"nearby",
		"none",
		"toggle",
		"visual",
	})
}
