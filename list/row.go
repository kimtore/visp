package list

type DataType string

const (
	// one for each data type
	DataTypeFIXME      DataType = "FIXME"
	DataTypeWindow              = "window"
	DataTypeList                = "list"
	DataTypeLogLine             = "logline"
	DataTypeKeyBinding          = "keybinding"
	DataTypeTrack               = "track"
	DataTypeDevice              = "device"
	DataTypeAlbum               = "album"
	DataTypePlaylist            = "playlist"
)

type Row interface {
	// Unique identifier of the data in the row.
	ID() string

	SetID(string)

	// Return the dataset.
	Fields() map[string]string

	SetFields(map[string]string)

	// Return the keys of the data.
	Keys() []string

	// Indicates what kind of data this row represents.
	Kind() DataType

	// Set a value in the dataset.
	Set(key, value string)

	// Get a value from the dataset.
	Get(key string) string
}

var _ Row = &BaseRow{}

type BaseRow struct {
	id     string
	kind   DataType
	fields map[string]string
}

func (row *BaseRow) Fields() map[string]string {
	return row.fields
}

func NewRow(id string, kind DataType, fields map[string]string) *BaseRow {
	if fields == nil {
		fields = make(map[string]string)
	}
	return &BaseRow{
		id:     id,
		kind:   kind,
		fields: fields,
	}
}

func (row *BaseRow) ID() string {
	return row.id
}

func (row *BaseRow) Kind() DataType {
	return row.kind
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
