package commands_test

import (
	"testing"

	"github.com/ambientsound/visp/commands"
)

var bindTests = []commands.Test{
	// Valid forms
	{`global foo bar`, true, nil, nil, []string{}},
	{`global foo bar baz`, true, nil, nil, []string{}},
	{`global []{}$|"test" foo bar`, true, nil, nil, []string{}},

	// Invalid forms
	{`x`, false, nil, nil, []string{}},
	{`global bar`, false, nil, nil, []string{}},
}

func TestBind(t *testing.T) {
	commands.TestVerb(t, "bind", bindTests)
}
