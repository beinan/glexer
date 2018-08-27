package glexer

type Schema struct {
	Types map[string]Type
}

type Type struct {
	Desc      string
	Name      string
	Category  string //type, enum, input, scalar, interface...
	Interface string
	Fields    map[string]Field
}

type Field struct {
	Desc string
	Name string
	Type FieldType
}

type FieldType struct {
	Name           string
	NonNull        bool
	IsArray        bool
	ElementNonNull bool
}

type Comment struct {
	Text string
}
