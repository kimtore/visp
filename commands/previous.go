package commands

import (
	"context"

	"github.com/ambientsound/visp/api"
)

// Previous instructs the player to go to the previous song.
type Previous struct {
	command
	api api.API
}

func NewPrevious(api api.API) Command {
	return &Previous{
		api: api,
	}
}

func (cmd *Previous) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Previous) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	return client.Previous(context.TODO())
}
