package input_test

import (
	"testing"

	"github.com/ambientsound/visp/api"
	"github.com/ambientsound/visp/input"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestCLISet tests that input.Interpreter registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestCLISet(t *testing.T) {
	var err error

	a := &api.MockAPI{}
	v := viper.New()
	a.On("Options").Return(v)
	a.On("Changed", api.ChangeOption, "foo").Return()

	opts := a.Options()
	iface := input.NewCLI(a)

	opts.Set("foo", "this string must die")

	err = iface.Exec("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.GetString("foo"))
}
