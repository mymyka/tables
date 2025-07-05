package schema

type Column struct {
	Name     string
	Type     string
	Nullable bool
}

type Table struct {
	Name    string
	Columns []Column
}
