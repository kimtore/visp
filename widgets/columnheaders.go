package widgets

import (
	"strings"

	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type ColumnheadersWidget struct {
	columns []list.Column
	view    views.View

	style.Styled
	views.WidgetWatchers
}

func NewColumnheadersWidget() (c *ColumnheadersWidget) {
	c = &ColumnheadersWidget{}
	c.columns = make(list.Columns, 0)
	return
}

func (c *ColumnheadersWidget) SetColumns(cols list.Columns) {
	c.columns = cols
}

func (c *ColumnheadersWidget) Draw() {
	x := 0
	y := 0
	for i := range c.columns {
		col := &c.columns[i]
		title := []rune(strings.Title(col.Tag))
		for p, r := range title {
			c.view.SetContent(x+p, y, r, nil, c.Style("header"))
		}
		x += col.Width()
	}
}

func (c *ColumnheadersWidget) SetView(v views.View) {
	c.view = v
}

func (c *ColumnheadersWidget) Size() (int, int) {
	x, y := c.view.Size()
	y = 1
	return x, y
}

func (w *ColumnheadersWidget) Resize() {
}

func (w *ColumnheadersWidget) HandleEvent(ev tcell.Event) bool {
	return false
}
