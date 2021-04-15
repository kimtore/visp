package db

import (
	"github.com/ambientsound/visp/list"
)

type Row struct {
	*list.BaseRow
	list list.List
}

func (row *Row) List() list.List {
	return row.list
}
