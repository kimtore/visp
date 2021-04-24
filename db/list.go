package db

import (
	"strconv"

	"github.com/ambientsound/visp/list"
)

type List struct {
	list.Base
	lists map[string]list.List
	last  list.List
}

var _ list.List = &List{}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("windows")
	this.SetName("Windows")
	this.SetVisibleColumns([]string{"name", "size"})
	this.lists = make(map[string]list.List)
	return this
}

func NewRow(lst list.List) list.Row {
	return &Row{
		BaseRow: list.NewRow(lst.ID(), map[string]string{
			"name": lst.Name(),
			"size": strconv.Itoa(lst.Len()),
		}),
		list: lst,
	}
}

// Cache adds a list to the database. Returns the row number of the list.
func (s *List) Cache(lst list.List) int {
	existing := s.RowByID(lst.ID())
	if existing == nil {
		s.Add(NewRow(lst))
		return s.Len() - 1
	} else {
		n, _ := s.RowNum(lst.ID())
		return n
	}
}

func (s *List) Current() list.List {
	row := s.CursorRow()
	if row == nil {
		return nil
	}
	return row.(*Row).List()
}

func (s *List) List(id string) list.List {
	row := s.RowByID(id)
	if row == nil {
		return nil
	}
	return row.(*Row).List()
}

func (s *List) SetLast(last list.List) {
	s.last = last
}

func (s *List) Last() list.List {
	return s.last
}
