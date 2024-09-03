// lexer/lexer_test.go

package lexer

import (
  "os"
	"testing"
  "github.com/ayase-mstk/go32as/src/parse"
)

/*
  コメントがUnicodeの場合があるので、rune単位で読み込む
*/
func  readFile(filename string) ([]rune, error) {
  data, err := os.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  runes := []rune(string(data))
  return runes, nil
}

type lexerTestStruct struct {
		expectedType  parse.TokenType
		expectedVal   string
    expectedRow   int
}

func expectSameSize(t *testing.T, got, want int) {
  if got != want {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      got, want)
  }
}

func expectSame(t *testing.T, got []parse.IToken, want []lexerTestStruct) {
  expectSameSize(t, len(got), len(want))

	for i, tt := range want {
    actual := got[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. got=%q, expected=%q",
        i, actual.Type(), tt.expectedType)
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. got=%q, expected=%q",
        i, actual.Val(), tt.expectedVal)
    }

    if actual.Row() != tt.expectedRow {
      t.Fatalf("test[%d] - row wrong. got=%d, expected=%d",
        i, actual.Row(), tt.expectedRow)
    }
  }
}


func TestLexer1(t *testing.T) {
	input := []rune("")
	tokens, err := parse.LexerLine(input, 1)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

  expectSameSize(t, len(tokens), 0)
}

func TestLexer2(t *testing.T) {
	input := []rune("main: addi a0, a1, 42")
	tokens, err := parse.LexerLine(input, 1)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType:   parse.LABEL,
      expectedVal: "main:",
      expectedRow:  1,
    },
		{
      expectedType:   parse.OPECODE,
      expectedVal: "addi",
      expectedRow:  1,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "a0",
      expectedRow:  1,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "a1",
      expectedRow:  1,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "42",
      expectedRow:  1,
    },
	}

  expectSame(t, tokens, tests)
}

func TestLexer3(t *testing.T) {
	input := []rune("lw x6, 0(x7)")
	tokens, err := parse.LexerLine(input, 2)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType:   parse.OPECODE,
      expectedVal: "lw",
      expectedRow:  2,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "x6",
      expectedRow:  2,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "0(x7)",
      expectedRow:  2,
    },
	}

  expectSame(t, tokens, tests)
}

func TestLexer4(t *testing.T) {
	input := []rune("var1: .word 100")
	tokens, err := parse.LexerLine(input, 1)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType: parse.LABEL,
      expectedVal:  "var1:",
      expectedRow:    1,
    },
		{
      expectedType:     parse.DIRECTIVE,
      expectedVal:  ".word",
      expectedRow:    1,
    },
		{
      expectedType:     parse.UNKNOWN,
      expectedVal:  "100",
      expectedRow:    1,
    },
	}

  expectSame(t, tokens, tests)
}

func TestLexer5(t *testing.T) {
	input := []rune("str_hello: .string \"Hello World\"")
	tokens, err := parse.LexerLine(input, 1)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType: parse.LABEL,
      expectedVal:  "str_hello:",
      expectedRow:  1,
    },
		{
      expectedType: parse.DIRECTIVE,
      expectedVal:  ".string",
      expectedRow:  1,
    },
		{
      expectedType: parse.UNKNOWN,
      expectedVal:  "\"Hello World\"",
      expectedRow:  1,
    },
	}

  expectSame(t, tokens, tests)
}

func TestLexer6(t *testing.T) {
	input := []rune(".LC0: .string   \"I'm from %d\n\" # this is comment\r\n")
	tokens, err := parse.LexerLine(input, 1)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType: parse.LABEL,
      expectedVal:  ".LC0:",
      expectedRow:  1,
    },
		{
      expectedType: parse.DIRECTIVE,
      expectedVal:  ".string",
      expectedRow:  1,
    },
    {
      expectedType: parse.UNKNOWN,
      expectedVal:  "\"I'm from %d\n\"",
      expectedRow:  1,
    },
  }

  expectSame(t, tokens, tests)
}

func TestLexer7(t *testing.T) {
	input := []rune(".string \"abc\"#this is comment")
	tokens, err := parse.LexerLine(input, 3)
  if err != nil {
    t.Fatalf("test - lexer failed:\n%q", err.Error())
  }

	tests := []lexerTestStruct{
		{
      expectedType:   parse.DIRECTIVE,
      expectedVal:    ".string",
      expectedRow:    3,
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal:    "\"abc\"",
      expectedRow:    3,
    },
	}

  expectSame(t, tokens, tests)
}
