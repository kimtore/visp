package commands

import (
	"context"
	"fmt"

	"github.com/ambientsound/visp/log"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input/lexer"
)

// https://github.com/zmb3/spotify/v2/issues/159
const (
	RepeatOff     = "off"
	RepeatContext = "context"
	RepeatTrack   = "track"
)

// Repeat sets the playback repeat mode.
type Repeat struct {
	command
	api    api.API
	action string
}

// NewRepeat returns Repeat.
func NewRepeat(api api.API) Command {
	return &Repeat{
		api: api,
	}
}

func (cmd *Repeat) nextState(state string) string {
	switch state {
	default:
		log.Debugf("Current repeat context is unknown '%s', defaulting to 'off'", state)
		fallthrough
	case RepeatOff:
		return RepeatContext
	case RepeatContext:
		return RepeatTrack
	case RepeatTrack:
		return RepeatOff
	}
}

// Parse implements Command.
func (cmd *Repeat) Parse() error {

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteAction(lit)

	switch tok {
	case lexer.TokenIdentifier:
		break
	case lexer.TokenEnd:
		cmd.action = cmd.nextState(cmd.api.PlayerStatus().RepeatState)
		return nil
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	switch lit {
	case RepeatContext:
		cmd.action = RepeatContext
	case RepeatOff:
		cmd.action = RepeatOff
	case RepeatTrack:
		cmd.action = RepeatTrack
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()
	return cmd.ParseEnd()

}

// Exec implements Command.
func (cmd *Repeat) Exec() error {

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	err = client.Repeat(context.TODO(), cmd.action)
	if err != nil {
		return err
	}

	switch cmd.action {
	case RepeatTrack:
		log.Infof("Repeating current track.")
	case RepeatContext:
		log.Infof("Repeating current play context.")
	case RepeatOff:
		log.Infof("Repeat turned off.")
	}

	return nil
}

// setTabCompleteAction sets the tab complete list to available actions.
func (cmd *Repeat) setTabCompleteAction(lit string) {
	list := []string{
		RepeatContext,
		RepeatOff,
		RepeatTrack,
	}
	cmd.setTabComplete(lit, list)
}
