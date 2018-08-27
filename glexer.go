package glexer

import (
	"fmt"
	"strings"
	"text/scanner"
)

func ParseSchema(schema string) error {

	lexer := InitLexer(schema)
	for lexer.Next() {
		schema := Try(ASchema)(lexer)
		fmt.Printf("\n types collected %#v", schema)
		//fmt.Println("At position", s.Pos(), ":", s.TokenText(), "--", scanner.TokenString(tok))
	}
	return nil
}

func ParseQuery(query string) error {
	s := &scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings,
	}
	s.Init(strings.NewReader(query))

	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()
		fmt.Println("At position", s.Pos(), ":", s.TokenText(), "--", scanner.TokenString(tok))
	}
	return nil
}
