package schema

type Table struct {
	Name string

	Fields  []*Field
	Indexes []*Index
	Primary []*Field
}

type Field struct {
	Name     string
	Type     FieldType
	Primary  bool
	Auto     bool
	Size     int
	Nullable bool
	Default  string
}

type Index struct {
	Name   string
	Unique bool

	Fields []*Field
}
