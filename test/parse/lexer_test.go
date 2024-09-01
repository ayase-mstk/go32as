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


func TestLexer1(t *testing.T) {
	input := []rune("")
	tokens := parse.LexerLine(input)

  if len(tokens) != 0 {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), 0)
  }
}

func TestLexer2(t *testing.T) {
	input := []rune("main: addi a0, a1, 42")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:   parse.LABEL,
      expectedVal: "main:",
    },
		{
      expectedType:   parse.OPECODE,
      expectedVal: "addi",
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "a0",
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "a1",
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "42",
    },
	}

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}

func TestLexer3(t *testing.T) {
	input := []rune("lw x6, 0(x7)")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:   parse.OPECODE,
      expectedVal: "lw",
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "x6",
    },
		{
      expectedType:   parse.UNKNOWN,
      expectedVal: "0(x7)",
    },
	}

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}

func TestLexer4(t *testing.T) {
	input := []rune("var1: .word 100")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:     parse.LABEL,
      expectedVal:  "var1:",
    },
		{
      expectedType:     parse.DIRECTIVE,
      expectedVal:  ".word",
    },
		{
      expectedType:     parse.UNKNOWN,
      expectedVal:  "100",
    },
	}

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}

func TestLexer5(t *testing.T) {
	input := []rune("str_hello: .string \"Hello World\"")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:     parse.LABEL,
      expectedVal:  "str_hello:",
    },
		{
      expectedType:     parse.DIRECTIVE,
      expectedVal:  ".string",
    },
		{
      expectedType:     parse.UNKNOWN,
      expectedVal:  "\"Hello World\"",
    },
	}

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}

func TestLexer6(t *testing.T) {
	input := []rune(".LC0: .string   \"I'm from %d\n\" # this is comment\r\n")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:     parse.LABEL,
      expectedVal:  ".LC0:",
    },
		{
      expectedType:     parse.DIRECTIVE,
      expectedVal:  ".string",
    },
    {
      expectedType:     parse.UNKNOWN,
      expectedVal:  "\"I'm from %d\n\"",
    },
  }

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}

func TestLexer7(t *testing.T) {
	input := []rune(".string \"abc\"#this is comment")
	tokens := parse.LexerLine(input)

	tests := []struct {
		expectedType    parse.TokenType
		expectedVal string
	}{
		{
      expectedType:     parse.DIRECTIVE,
      expectedVal:  ".string",
    },
		{
      expectedType:     parse.UNKNOWN,
      expectedVal:  "\"abc\"",
    },
	}

  if len(tokens) != len(tests) {
    t.Fatalf("token num wrong. actual = %d, expected = %d",
      len(tokens), len(tests))
  }

	for i, tt := range tests {
    actual := tokens[i]

    if actual.Type() != tt.expectedType {
      t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q",
        i, tt.expectedType, actual.Type())
    }

    if actual.Val() != tt.expectedVal {
      t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
        i, tt.expectedVal, actual.Val())
    }
	}
}
