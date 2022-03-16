package commands

import (
	"context"

	"github.com/ambientsound/visp/api"
)

// Pause toggles play/paused state.
type Pause struct {
	command
	api api.API
}

func NewPause(api api.API) Command {
	return &Pause{
		api: api,
	}
}

func (cmd *Pause) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Pause) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	if cmd.api.PlayerStatus().Playing {
		return client.Pause(context.TODO())
	} else {
		return client.Play(context.TODO())
	}
}
