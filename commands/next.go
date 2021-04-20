package commands

import (
	"github.com/ambientsound/visp/api"
)

// Next instructs the player to go to the next song.
type Next struct {
	command
	api api.API
}

func NewNext(api api.API) Command {
	return &Next{
		api: api,
	}
}

func (cmd *Next) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Next) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	defer cmd.api.Changed(api.ChangePlayerStateInvalid, nil)

	return client.Next()
}
