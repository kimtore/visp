package keys

import (
	"github.com/ambientsound/pms/list"
)

func (s *Sequencer) List() list.List {
	this := list.New()

	this.Clear()
	this.SetID("keybindings")
	this.SetName("Key bindings")
	this.SetVisibleColumns([]string{"context", "keySequence", "command"})

	for _, bind := range s.binds {
		this.Add(list.Row{
			"context":     bind.Context,
			"keySequence": bind.Sequence.String(),
			"command":     bind.Command,
		})
	}

	this.Sort([]string{"keySequence", "context"})

	return this
}
