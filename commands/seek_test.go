package commands_test

import (
	"testing"

	"github.com/ambientsound/visp/commands"
	"github.com/ambientsound/visp/player"
)

var seekTests = []commands.Test{
	// Valid forms
	{`-2`, true, setupTestSeek, nil, []string{}},
	{`+13`, true, setupTestSeek, nil, []string{}},
	{`1329`, true, setupTestSeek, nil, []string{}},

	// Invalid forms
	{`nan`, false, setupTestSeek, nil, []string{}},
	{`+++1`, false, setupTestSeek, nil, []string{}},
	{`-foo`, false, setupTestSeek, nil, []string{}},
	{`$1`, false, setupTestSeek, nil, []string{}},
}

func setupTestSeek(data *commands.TestData) {
	data.MockAPI.On("PlayerStatus").Return(player.State{})
}

func TestSeek(t *testing.T) {
	commands.TestVerb(t, "seek", seekTests)
}
