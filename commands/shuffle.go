package commands

import (
	"context"
	"fmt"

	"github.com/ambientsound/visp/log"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// Shuffle sets the playback shuffle mode.
type Shuffle struct {
	command
	api    api.API
	action bool
}

// NewShuffle returns Shuffle.
func NewShuffle(api api.API) Command {
	return &Shuffle{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Shuffle) Parse() error {

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteAction(lit)

	switch tok {
	case lexer.TokenIdentifier:
		break
	case lexer.TokenEnd:
		cmd.action = !cmd.api.PlayerStatus().ShuffleState
		return nil
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	switch lit {
	case "on":
		cmd.action = true
	case "off":
		cmd.action = false
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()
	return cmd.ParseEnd()

}

// Exec implements Command.
func (cmd *Shuffle) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	err = client.Shuffle(context.TODO(), cmd.action)
	if err != nil {
		return err
	}

	if cmd.action {
		log.Infof("Shuffle turned on.")
	} else {
		log.Infof("Shuffle turned off.")
	}

	return nil
}

// setTabCompleteAction sets the tab complete list to available actions.
func (cmd *Shuffle) setTabCompleteAction(lit string) {
	list := []string{
		"off",
		"on",
	}
	cmd.setTabComplete(lit, list)
}
