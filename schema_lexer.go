package glexer

import (
	"fmt"
)

func ASchema(l *Lexer) Result {
	types := Many(AType)(l).Data.([]Result)
	typesMap := make(map[string]Type)
	for _, typeResult := range types {
		typeItem := typeResult.Data.(Type)
		typesMap[typeItem.Name] = typeItem
	}
	return ResultFound(
		Schema{
			Types: typesMap,
		},
		l.Pos(),
	)
}

func collectMultiLineComment(l *Lexer) string {
	comments := Many(AComment)(l).Data.([]Result)
	fmt.Println("collect comments", comments)
	desc := ""
	for _, comment := range comments {
		desc += comment.Data.(Comment).Text + "\n"
	}
	return desc
}

func collectFields(l *Lexer) map[string]Field {
	fieldsResult := Enforce(AFieldBlock, "Fields are required for type")(l)
	return fieldsResult.Data.(map[string]Field)
}

func collectTypeName(l *Lexer) string {
	r := Enforce(AnIdent, "Type name is required")(l)
	return r.Data.(string)
}
func AType(l *Lexer) Result {
	desc := collectMultiLineComment(l)
	var typename string
	var fields map[string]Field
	category := l.Text //type category such as type, enum, input and etc..
	switch category {
	case "schema":
		typename = "schema" //schema root, just a name as placeholder
	case "enum":
		l.Next()
		typename = collectTypeName(l)
		//todo: collect enum element
	case "scalar":
		l.Next()
		typename = collectTypeName(l)
	case "type":
		l.Next()
		typename = collectTypeName(l)
		//todo inheritance
		fields = collectFields(l)
	case "input":
		fallthrough
	case "interface":
		l.Next()
		typename = collectTypeName(l)
	default:
		//for other starting keyword, it's not a "type"
		return ResultNotFound
	}
	fmt.Println("a type name:", typename)

	fmt.Printf("fields: %+v", fields)
	return ResultFound(
		Type{
			Desc:     desc,
			Name:     typename,
			Fields:   fields,
			Category: category,
		},
		l.Pos(),
	)
}

func AFieldBlock(l *Lexer) Result {
	if l.Tok != '{' {
		return ResultNotFound
	}
	l.Next()
	fields := Many(AField)(l).Data.([]Result)
	fieldsMap := make(map[string]Field)
	for _, fieldResult := range fields {
		field := fieldResult.Data.(Field)
		fieldsMap[field.Name] = field
	}
	Enforce(ARune('}'), "Expect a '}' here")(l)
	return ResultFound(fieldsMap, l.Pos())
}

func AField(l *Lexer) Result {
	if !l.IsIdent() {
		return ResultNotFound
	}
	fieldName := l.TextAndMove()
	Enforce(ARune(':'), "Expect a ':' here")(l)
	fieldType := Enforce(AFieldType, "Expect a type for the field")(l)
	return ResultFound(Field{
		Name: fieldName,
		Type: fieldType.Data.(FieldType),
	}, l.Pos())
}

func isNonNull(l *Lexer) bool {
	var isNonNull bool
	if l.Tok == '!' {
		isNonNull = true
		l.Next()
	}
	return isNonNull
}
func AFieldType(l *Lexer) Result {
	if l.Tok == '[' { //array type
		l.Next()
		elementType := Enforce(AFieldType, "Expect a type for the array element")(l)
		isNonNull := isNonNull(l)
		Enforce(ARune(']'), "Expect a ]")(l)
		et := elementType.Data.(FieldType)
		return ResultFound(
			FieldType{
				Name:           et.Name,
				NonNull:        isNonNull,
				IsArray:        true,
				ElementNonNull: et.NonNull,
			},
			l.Pos(),
		)
	} else { //non-array type
		typeName := collectTypeName(l)
		isNonNull := isNonNull(l)
		return ResultFound(FieldType{
			Name:    typeName,
			NonNull: isNonNull,
		},
			l.Pos(),
		)
	}
}

func AComment(l *Lexer) Result {
	fmt.Println("comment begin:", l.Text)
	if l.Text != "#" {
		return ResultNotFound
	}
	comment := l.LineStr()
	fmt.Println("comment", comment)
	return ResultFound(
		Comment{
			Text: comment,
		},
		l.Pos(),
	)
}
