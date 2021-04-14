package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var selectTests = []commands.Test{
	// Valid forms
	{`visual`, true, nil, nil, []string{}},
	{`toggle`, true, nil, nil, []string{}},
	{`nearby artist tit`, true, initSongTags, nil, []string{"title"}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`visual 1`, false, nil, nil, []string{}},
	{`toggle 1`, false, nil, nil, []string{}},
	{`nearby`, false, nil, nil, []string{}},
}

func TestSelect(t *testing.T) {
	commands.TestVerb(t, "select", selectTests)
}
