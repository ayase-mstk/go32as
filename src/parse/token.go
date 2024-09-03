package parse

import "fmt"

type TokenType int

const (
  LABEL TokenType = iota
  DIRECTIVE
  OPECODE
  OPERAND
  UNKNOWN
)

type ArgFlag int

const (
  Take0  ArgFlag = 1 << iota // 00000001
  Take1                      // 00000010
  Take2                      // 00000100
  Take3                      // 00001000
  TakeN                      // 00010000
  Take12                     // 00100000
  Take123                    // 01000000
)

type IToken interface {
  Type()    TokenType
  Val()     string
  Row()     int
}

type Token struct {
  typ    TokenType
  val    string
  row    int
}

func (t Token) Type() TokenType { return t.typ }
func (t Token) Val() string { return t.val }
func (t Token) Row() int { return t.row }

func newToken(ttype TokenType, val string, row int) Token {
  return Token{typ: ttype, val: val, row: row}
}

func printTokens(tokens []Token) {
  for _, t := range tokens {
    t.printFeature()
  }
}

func (t *Token) printFeature() {
  fmt.Println("val=", t.Val())
  fmt.Println("type=", t.Type())
  fmt.Println("row=", t.Row())
}
