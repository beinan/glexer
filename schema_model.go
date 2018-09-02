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
	Args map[string]Argument
}

type FieldType struct {
	Name           string
	NonNull        bool
	IsArray        bool
	ElementNonNull bool
}

type Argument struct {
	Name         string
	Type         FieldType
	DefaultValue string //todo: using graphql_model types
}

type Comment struct {
	Text string
}
