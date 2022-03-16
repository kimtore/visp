package db

import (
	"sort"
	"strconv"

	"github.com/ambientsound/visp/list"
	spotify_playlists "github.com/ambientsound/visp/spotify/playlists"
	"github.com/zmb3/spotify/v2"
)

type List struct {
	list.Base
	lists      map[string]list.List
	last       list.List
	nameLookup map[string]interface{}
}

var _ list.List = &List{}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("windows")
	this.SetName("Windows")
	this.SetVisibleColumns([]string{"name", "size"})
	this.nameLookup = make(map[string]interface{})
	this.lists = make(map[string]list.List)
	return this
}

func NewRow(lst list.List) list.Row {
	return &Row{
		BaseRow: list.NewRow(lst.ID(), list.DataTypeWindow, map[string]string{
			"name": lst.Name(),
			"size": strconv.Itoa(lst.Len()),
		}),
		list: lst,
	}
}

// Cache adds a list to the database. Returns the row number of the list.
func (s *List) Cache(lst list.List) int {
	defer func() {
		s.addNameLookupChildren(lst)
		s.nameLookup[lst.Name()] = lst
		// log.Debugf("List lookup table: %v", s.nameLookup)
	}()
	existing := s.RowByID(lst.ID())
	if existing == nil {
		s.Add(NewRow(lst))
		return s.Len() - 1
	} else {
		ex := existing.(*Row)
		n, _ := s.RowNum(lst.ID())
		row := s.Row(n)
		delete(s.nameLookup, row.Get("name"))
		ex.list = lst
		row.Set("name", lst.Name())
		row.Set("size", strconv.Itoa(lst.Len()))
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

func (s *List) Names() []string {
	names := make([]string, 0, len(s.nameLookup))
	for k := range s.nameLookup {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func (s *List) Lookup(name string) interface{} {
	return s.nameLookup[name]
}

func (s *List) addNameLookupChildren(lst list.List) {
	playlists, ok := lst.(*spotify_playlists.List)
	if !ok {
		return
	}
	for _, row := range playlists.All() {
		s.nameLookup[row.Get("name")] = spotify.ID(row.ID())
	}
}
