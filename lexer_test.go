package glexer

import "testing"

func TestParseSchema(t *testing.T) {

}

func TestParseQuery(t *testing.T) {
	err := ParseQuery("")
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}
