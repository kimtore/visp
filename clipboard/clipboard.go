package clipboard

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ambientsound/visp/list"
)

type List struct {
	list.Base
	lists  map[string]list.List
	active list.List
}

var _ list.List = &List{}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("clipboards")
	this.SetName("Clipboards")
	this.SetVisibleColumns([]string{"name", "size", "time"})
	this.lists = make(map[string]list.List)
	return this
}

func (clipboard *List) Insert(l list.List) {
	index := strconv.Itoa(clipboard.Len() + 1)

	l.SetName(fmt.Sprintf("Clipboard %s", index))
	l.SetID(fmt.Sprintf("clipboard_%s", index))

	clipboard.Add(list.NewRow(
		index,
		map[string]string{
			"name": l.Name(),
			"size": strconv.Itoa(l.Len()),
			"time": time.Now().Format(time.RFC1123),
		},
	))

	clipboard.lists[index] = l

	clipboard.active = l
}

func (clipboard *List) Update(l list.List) {
	row := clipboard.RowByID(l.ID())
	if row == nil {
		return
	}
	row.Set("name", l.Name())
	row.Set("size", strconv.Itoa(l.Len()))
}

// Return the currently selected clipboard
func (clipboard *List) Current() list.List {
	row := clipboard.CursorRow()
	if row == nil {
		return nil
	}
	return clipboard.lists[row.ID()]
}

func (clipboard *List) Get(id string) list.List {
	return clipboard.lists[id]
}

func (clipboard *List) Active() list.List {
	return clipboard.active
}
