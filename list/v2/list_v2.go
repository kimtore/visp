package list_v2

type Cell interface {
	String()
}

type Row struct {
	cells map[string]Cell
}

type Grid struct {
	id      string
	headers Headers
	data    []Row
	source  interface{}
}

type Source interface {
	Rows() []Row
}

func New(source ...Source) *Grid {
	grid := &Grid{
		data: make([]Row, 0),
	}
	for _, s := range source {
		grid.data = append(grid.data, s.Rows()...)
	}
	return grid
}
