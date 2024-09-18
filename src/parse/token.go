package parse

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TLabel TokenType = iota
	TDirective
	TOpecode
	TUnknown
)

type Token struct {
	typ TokenType
	val string
	//row    int
}

func newToken(val string) Token {
	tk := Token{val: val}
	tk.setType()
	return tk
}

func (t Token) Type() TokenType { return t.typ }

func (t Token) Val() string { return t.val }

//func (t Token) Row() int { return t.row }

func (t Token) isLabel() bool {
	literal := t.Val()
	// ラベルかどうか
	if !strings.HasSuffix(literal, ":") {
		return false
	}
	literal = literal[:len(literal)-1]
	// すべて数値
	if isNumericStr(literal[:len(literal)-1]) {
		return true
	}
	// 接頭辞はalphabetかアンダーバー
	if isAlpha(literal[0]) || literal[0] == '_' || literal[0] == '.' {
		literal = literal[1:]
		for i := 0; i < len(literal); i++ {
			if !(isAlpha(literal[i]) || isNumeric(literal[i]) || literal[i] == '_' || literal[i] == '.') {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func (t Token) isDirective() bool {
	_, exists := directiveSet[t.Val()]
	return exists
}

func (t Token) isOpecode() bool {
	_, exists := OpecodeMap[t.Val()]
	return exists
}

func (t *Token) setType() {
	if t.isLabel() {
		t.typ = TLabel
	} else if t.isDirective() {
		t.typ = TDirective
	} else if t.isOpecode() {
		t.typ = TOpecode
	} else {
		t.typ = TUnknown
	}
}

func printTokens(tokens []Token) {
	for _, t := range tokens {
		printToken(t)
	}
}
func printToken(tk Token) {
	fmt.Printf("val=[%s]", tk.Val())
	fmt.Println("type=", tk.Type())
}
