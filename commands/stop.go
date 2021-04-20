package commands

import (
	"github.com/ambientsound/visp/api"
)

// Stop stops song playback.
type Stop struct {
	command
	api api.API
}

// NewStop returns Stop.
func NewStop(api api.API) Command {
	return &Stop{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Stop) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Stop) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	return client.Pause()
}
