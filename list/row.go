package list

type RowType int

const (
	RowTypePlain RowType = iota
)

type Row interface {
	ID() string
	SetID(string)
	Fields() map[string]string
	SetFields(map[string]string)
	Keys() []string
	Set(key, value string)
	Get(key string) string
}

var _ Row = &BaseRow{}

type BaseRow struct {
	id     string
	fields map[string]string
}

func (row *BaseRow) Fields() map[string]string {
	return row.fields
}

func NewRow(id string, fields map[string]string) *BaseRow {
	if fields == nil {
		fields = make(map[string]string)
	}
	return &BaseRow{
		id:     id,
		fields: fields,
	}
}

func (row *BaseRow) ID() string {
	return row.id
}

func (row *BaseRow) SetID(id string) {
	row.id = id
}

func (row *BaseRow) Set(key, value string) {
	row.fields[key] = value
}

func (row *BaseRow) Get(key string) string {
	return row.fields[key]
}

func (row *BaseRow) SetFields(fields map[string]string) {
	row.fields = fields
}

func (row *BaseRow) Keys() []string {
	keys := make([]string, 0)
	for k := range row.fields {
		keys = append(keys, k)
	}
	return keys
}
