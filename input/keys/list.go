package keys

import (
	"github.com/ambientsound/visp/list"
)

func (s *Sequencer) List() list.List {
	this := list.New()

	this.Clear()
	this.SetID("keybindings")
	this.SetName("Key bindings")
	this.SetVisibleColumns([]string{"context", "keySequence", "command"})

	for _, bind := range s.binds {
		keySequence := bind.Sequence.String()
		this.Add(list.NewRow(keySequence, map[string]string{
			"context":     bind.Context,
			"keySequence": keySequence,
			"command":     bind.Command,
		}))
	}

	this.Sort([]string{"keySequence", "context"})

	return this
}
