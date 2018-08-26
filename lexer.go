package glexer

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Lexer struct {
	src  string
	s    *scanner.Scanner
	Tok  rune   //last rune
	Text string //the token text scanner just scanned
}

func InitLexer(text string) *Lexer {
	s := &scanner.Scanner{
		Mode: scanner.ScanIdents,
	}
	s.Init(strings.NewReader(text))
	return &Lexer{
		src: text,
		s:   s,
	}
}

func (l *Lexer) Next() bool {
	l.Tok = l.s.Scan()
	if l.Tok == scanner.EOF {
		return false
	}
	l.Text = l.s.TokenText()
	return true
}

func (l *Lexer) IsIdent() bool {
	return l.Tok == scanner.Ident
}

func (l *Lexer) TextAndMove() string {
	text := l.Text
	l.Next()
	return text
}

func (l *Lexer) LineStr() string {
	line := ""
	for {
		next := l.s.Next()
		if next == '\n' || next == scanner.EOF {
			l.Next() //go to next line
			break
		}
		line += string(next)
	}
	return line
}

func (l *Lexer) Pos() Position {
	return Position{
		Offset: l.s.Offset,
		Line:   l.s.Line,
		Column: l.s.Column,
	}
}

type Result struct {
	Found bool
	Data  interface{}
	Pos   Position
}

var ResultNotFound = Result{Found: false}

func ResultFound(data interface{}, pos Position) Result {
	return Result{
		Found: true,
		Data:  data,
		Pos:   pos,
	}
}

type Position struct {
	Offset int
	Line   int
	Column int
}

type TypeDef struct {
	Desc string
	Type interface{}
}

type CustomeType struct {
	Name   string
	Fields map[string]Field
}

type Field struct {
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

func ATypeDef(l *Lexer) Result {
	comments := Many(AComment)(l).Data.([]Result)
	desc := ""
	for _, comment := range comments {
		desc += comment.Data.(Comment).Text + "\n"
	}
	aType := ACustomType(l)
	if !aType.Found {
		return ResultNotFound
	}
	return ResultFound(
		TypeDef{
			Desc: desc,
			Type: aType.Data,
		},
		l.Pos(),
	)
}

func AnIdent(l *Lexer) Result {
	fmt.Println("An Ident curr token:", l.Text, scanner.TokenString(l.Tok))
	if l.Tok != scanner.Ident {
		return ResultNotFound
	}
	defer l.Next() //eat this Ident and then move next
	return ResultFound(l.Text, l.Pos())
}

func ACustomType(l *Lexer) Result {
	if l.Text != "type" {
		return ResultNotFound
	}
	fmt.Println("custom type:", l.Text)
	l.Next() //go to type name
	fmt.Println("custom type type name:", l.Text)

	//type name has to be an identifier
	typeName := Enforce(l, "Type name is required", AnIdent)

	//TODO: enheritance

	fields := Enforce(l, "Fields are required for type", AFieldBlock)
	fmt.Printf("fields: %+v", fields)
	return ResultFound(
		CustomeType{
			Name:   typeName.Data.(string),
			Fields: fields.Data.(map[string]Field),
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
	Enforce(l, "Expect a '}' here", ARune('}'))
	return ResultFound(fieldsMap, l.Pos())
}

func AField(l *Lexer) Result {
	if !l.IsIdent() {
		return ResultNotFound
	}
	fieldName := l.TextAndMove()
	Enforce(l, "Expect a ':' here", ARune(':'))
	fieldType := Enforce(l, "Expect a type for the field", AFieldType)
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
		elementType := Enforce(l, "Expect a type for the array element", AFieldType)
		isNonNull := isNonNull(l)
		Enforce(l, "Expect a ]", ARune(']'))
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
		typeName := Enforce(l, "Expect a valid type name.", AnIdent)
		isNonNull := isNonNull(l)
		return ResultFound(FieldType{
			Name:    typeName.Data.(string),
			NonNull: isNonNull,
		},
			l.Pos(),
		)
	}
}

func ARune(ch rune) func(*Lexer) Result {
	return func(l *Lexer) Result {
		fmt.Println("expect:", string(ch), ch, "actual", string(l.Tok), l.Tok)
		if l.Tok != ch {
			return ResultNotFound
		}
		l.Next()                        //eat this rune and move next
		return ResultFound(ch, l.Pos()) //we don't really need a result
	}
}

func AComment(l *Lexer) Result {
	if l.Text != "#" {
		return ResultNotFound
	}
	comment := l.LineStr()
	return Result{
		Data: Comment{
			Text: comment,
		},
	}
}

func Enforce(l *Lexer, panicMsg string, aFunc func(*Lexer) Result) Result {
	result := aFunc(l)
	if !result.Found {
		panic(panicMsg)
	}
	return result
}

func Many(aFunc func(*Lexer) Result) func(*Lexer) Result {
	return func(l *Lexer) Result {
		var results []Result
		for {
			aResult := aFunc(l)
			if !aResult.Found {
				break
			}
			results = append(results, aResult)
		}
		return ResultFound(results, l.Pos())
	}
}

func Try(aFunc func(*Lexer) Result) func(*Lexer) Result {
	return func(l *Lexer) (ret Result) {
		defer func() {
			if err := recover(); err != nil {
				for i, line := range strings.Split(l.src, "\n") {
					fmt.Printf("%v : %v \n", i, line)
					if i == l.Pos().Line-1 {
						fmt.Printf("Sytax Parsing Error: %v ^^^^ \n", err)
					}
				}
				ret = ResultNotFound
			}
		}()
		return aFunc(l)
	}
}
