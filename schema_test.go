package glexer

import (
	"testing"

	. "github.com/beinan/extractor"
)

func TestExDesc(t *testing.T) {
	var l *Lexer

	//one line desc
	l = InitLexer("#abc abcd\n cced\n")
	var result string
	if exDesc(&result)(l) == false || result != "abc abcd" {
		t.Fatalf("One line: ExDesc(result) should get abc abcd, but actual: %v", result)
	}

	result = ""
	//no desc
	l = InitLexer("abc dd")
	if exDesc(&result)(l) == false || result != "" {
		t.Fatalf("No line: ExDesc(result) should return empty string, actual %v", result)
	}

	//multi line desc
	l = InitLexer("# abc dd \n # aa")
	if exDesc(&result)(l) == false || result != "abc dd aa" {
		t.Fatalf("Multi line: ExDesc(result) should return abc dd aa, actual %v", result)
	}
}

func TestExScalar(t *testing.T) {
	var l *Lexer

	l = InitLexer("scalar AAA")

	var tp Type
	if exScalar(&tp)(l) == false || tp.Name != "AAA" || tp.Category != "scalar" {
		t.Fatalf("Incorrect scalar type extracted: %#v", tp)
	}
}

func TestExCustomType(t *testing.T) {
	var l *Lexer

	l = InitLexer("type AA {}")

	var tp Type
	if exCustomType(&tp)(l) == false || tp.Name != "AA" || tp.Category != "type" {
		t.Fatalf("Incorrect custom type extracted: %#v", tp)
	}

	l = InitLexer("type AA2 implements BB { a:Int b(c: Int = 100):[Int!]!}")

	if exCustomType(&tp)(l) == false || tp.Name != "AA2" ||
		tp.Interface != "BB" || tp.Category != "type" {
		t.Fatalf("Incorrect custom type extracted: %#v", tp)
	}
}
