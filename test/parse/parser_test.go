// parse/parser_test.go

package parsetest

import (
	"fmt"
	"github.com/ayase-mstk/go32as/src/parse"
	"testing"
)

type parseTestStruct struct {
	expectedType  parse.StmtType
	expectedLabel string
	expectedRow   int
}

func expectSame(t *testing.T, got []parse.Stmt, want []parseTestStruct) {
	for i, tt := range want {
		actual := got[i]

		if actual.Type() != tt.expectedType {
			t.Fatalf("test[%d] - StmtType wrong. got=%q, expected=%q",
				i, actual.Type(), tt.expectedType)
		}

		if actual.LSymbol() != tt.expectedLabel {
			t.Fatalf("test[%d] - label wrong. got=%q, expected=%q",
				i, actual.LSymbol(), tt.expectedLabel)
		}

		if actual.Row() != tt.expectedRow {
			t.Fatalf("test[%d] - row wrong. got=%d, expected=%d",
				i, actual.Row(), tt.expectedRow)
		}
	}
}

func TestParser1(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("main:")
	stmt, err := parse.ParseLine(input, 1)
	stmts = append(stmts, stmt)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseTestStruct{
		{
			expectedType:  parse.UNKNOWN,
			expectedLabel: "main:",
			expectedRow:   1,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser2(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("main:")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}
	input2 := []rune("    sw ra, 32(sp)")
	stmt2, err2 := parse.ParseLine(input2, 2)
	if err2 != nil {
		t.Fatalf("test - parse failed:\n%q", err2.Error())
	}
	stmts = append(stmts, stmt, stmt2)

	tests := []parseTestStruct{
		{
			expectedType:  parse.UNKNOWN,
			expectedLabel: "main:",
			expectedRow:   1,
		},
		{
			expectedType:  parse.OPERATION,
			expectedLabel: "",
			expectedRow:   2,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser3(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune(".L0:")
	stmt, err := parse.ParseLine(input, 11)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}
	input2 := []rune("    .section .text")
	stmt2, err2 := parse.ParseLine(input2, 12)
	if err2 != nil {
		t.Fatalf("test - parse failed:\n%q", err2.Error())
	}
	stmts = append(stmts, stmt, stmt2)

	tests := []parseTestStruct{
		{
			expectedType:  parse.UNKNOWN,
			expectedLabel: ".L0:",
			expectedRow:   11,
		},
		{
			expectedType:  parse.DIRECTIVE,
			expectedLabel: "",
			expectedRow:   12,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser4(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("var1: .word 100")
	stmt, err := parse.ParseLine(input, 1)
	stmts = append(stmts, stmt)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}

	tests := []parseTestStruct{
		{
			expectedType:  parse.DIRECTIVE,
			expectedLabel: "var1:",
			expectedRow:   1,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser5(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("str_hello: .string \"Hello World\"")
	stmt, err := parse.ParseLine(input, 1)
	stmts = append(stmts, stmt)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}

	tests := []parseTestStruct{
		{
			expectedType:  parse.DIRECTIVE,
			expectedLabel: "str_hello:",
			expectedRow:   1,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser6(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune(".LC0: .string   \"I'm from %d\n\" # this is comment\r\n")
	stmt, err := parse.ParseLine(input, 1)
	stmts = append(stmts, stmt)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}

	tests := []parseTestStruct{
		{
			expectedType:  parse.DIRECTIVE,
			expectedLabel: ".LC0:",
			expectedRow:   1,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser7(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune(".string \"abc\"#this is comment")
	stmt, err := parse.ParseLine(input, 3)
	stmts = append(stmts, stmt)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}

	tests := []parseTestStruct{
		{
			expectedType:  parse.DIRECTIVE,
			expectedLabel: "",
			expectedRow:   3,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser8(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("#this is comment")
	stmt, err := parse.ParseLine(input, 3)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}
	stmts = append(stmts, stmt)

	tests := []parseTestStruct{
		{
			expectedType:  parse.UNKNOWN,
			expectedLabel: "",
			expectedRow:   3,
		},
	}

	expectSame(t, stmts, tests)
}

func TestParser9(t *testing.T) {
	var stmts []parse.Stmt
	input := []rune("     lw x5, 10(x5)     #this is comment")
	stmt, err := parse.ParseLine(input, 3)
	if err != nil {
		t.Fatalf("test - parser failed:\n%q", err.Error())
	}
	stmts = append(stmts, stmt)

	tests := []parseTestStruct{
		{
			expectedType:  parse.OPERATION,
			expectedLabel: "",
			expectedRow:   3,
		},
	}

	expectSame(t, stmts, tests)
}



/*
==========================================
=============  Error Test  ===============
==========================================
*/

const (
	UnrecognizedError = "junk at end of line, first unrecognized character is `%c'"
)

func expectErrorMessage(t *testing.T, actual, expect string) {
	if actual != expect {
		t.Fatalf("test - lexer error msg is different from expected.\nactual: %q\nexpected: %q", actual, expect)
	}
}

func TestParserError1(t *testing.T) {
	input := []rune("main: .L0:")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parser didnot fail: expect fail")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '.'))
}

func TestParserError2(t *testing.T) {
	input := []rune("main: lw x6, 0(x7) .section")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parser didnot fail: expect fail")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '.'))
}

func TestParserError3(t *testing.T) {
	input := []rune("lw x6, 0(x7) main:")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parser didnot fail: expect fail")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'm'))
}

func TestParserError4(t *testing.T) {
	input := []rune(".text .L0:")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parser didnot fail: expect fail")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '.'))
}
