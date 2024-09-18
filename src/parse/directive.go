package parse

import (
	"errors"
	"fmt"
)

type Directive struct {
	name    string
	args    []string
	argTyps []DirectiveArgType
	src     []rune
	idx     int
	// 引数自体がvalidかどうかはparseで判断
}

func (d *Directive) Name() string                { return d.name }
func (d *Directive) Args() []string              { return d.args }
func (d *Directive) ArgTyps() []DirectiveArgType { return d.argTyps }
func (d *Directive) isEOF() bool {
	return d.idx == len(d.src)
}
func (d *Directive) nextVal() (string, DirectiveArgType) {
	isLiteral := false
	start := d.idx
	for ; d.idx < len(d.src); d.idx++ {
		if '"' == d.src[d.idx] && !isLiteral {
			isLiteral = true
		} else if '"' == d.src[d.idx] && isLiteral {
			isLiteral = false
		} else if isLiteral {
			continue
		}
		c := d.src[d.idx]
		if c == ' ' || c == '\t' || c == ',' || c == '(' || c == ')' {
			//if isDelim(o.src[o.idx]) {
			break
		}
	}
	// この関数に入った場合かならずtokenがある
	val := string(d.src[start:d.idx])
	typ := analyzeDirArgType(val)
	return val, typ
}
func (d *Directive) skipUntilNextVal() {
	for ; d.idx < len(d.src); d.idx++ {
		c := d.src[d.idx]
		if c == '#' { // comment以降は飛ばす
			d.idx = len(d.src)
			return
		} else if c != ' ' && c != '\t' && c != ',' && c != '(' && c != ')' {
			return
		}
	}
}
func analyzeDirArgType(val string) DirectiveArgType {
	if isImmediate(val) {
		return INT
	}
	return STR
}

func (d Directive) isSection() bool {
	return d.name == Text || d.name == Data || d.name == RoData || d.name == Bss
}

// directive
// map[string][]string
const (
	Align = ".align"
	// P2Align   = ".p2align"
	// BAlign    = ".balign"
	File    = ".file"
	Globl   = ".globl"
	Local   = ".local"
	Comm    = ".comm"
	Common  = ".common"
	Ident   = ".ident"
	Section = ".section"
	Size    = ".size"
	Text    = ".text"
	Data    = ".data"
	RoData  = ".rodata"
	Bss     = ".bss"
	String  = ".string"
	Asciz   = ".asciz"
	Equ     = ".equ"
	Macro   = ".macro"
	Endm    = ".endm"
	Type    = ".type"
	// Option     = ".option"
	Byte  = ".byte"
	Byte2 = ".2byte"
	Half  = ".half"
	Short = ".short"
	Byte4 = ".4byte"
	Word  = ".word"
	Long  = ".long"
	// Float      = ".float"
	// DtprelWord = ".dtprelword"
	Zero = ".zero"
	// VariantCC = ".variant_cc"
	Attribute = ".attribute"
)

type DirectiveArgType int

const (
	INT DirectiveArgType = 1 << iota
	STR
)

var directiveSet = map[string][]DirectiveArgType{
	Align: {INT},
	// P2Align:    {},
	// BAlign:     {},
	File:    {STR},
	Globl:   {STR},
	Local:   {STR},
	Comm:    {STR, INT, INT},
	Common:  {STR, INT, INT},
	Ident:   {STR},
	Section: {STR},      // STRというより、セクション名
	Size:    {STR, INT}, // とりあえずアドレス計算は対応しない
	Text:    {},
	Data:    {},
	RoData:  {},
	Bss:     {},
	String:  {STR},
	Asciz:   {STR},
	Equ:     {STR, INT},
	Macro:   {STR},
	Endm:    {},
	Type:    {STR, INT},
	// Option:     {},
	Byte:  {INT},
	Byte2: {INT},
	Half:  {INT},
	Short: {INT},
	Byte4: {INT},
	Word:  {INT},
	Long:  {INT},
	// Float:      {},
	// DtprelWord: {},
	Zero: {INT},
	// VariantCC: {STR},
	Attribute: {STR, INT},
}

var directiveMap = map[string]func() error{
	//Align:     parseAlign,
	//File:      parseFile,
	//Globl:     parseGlobl,
	//Local:     parseLocal,
	//Comm:      parseCommon,
	//Common:    parseCommon,
	//Ident:     parseIdent,
	//Section:   parseSection,
	//Size:      parseSize,
	//Text:      parseText,
	//Data:      parseData,
	//RoData:    parseRoData,
	//Bss:       parseBss,
	//String:    parseString,
	//Asciz:     parseAsciz,
	//Equ:       parseEqu,
	//Macro:     parseMacro,
	//Endm:      parseEndm,
	//Type:      parseType,
	//Byte:      parseByte,
	//Byte2:     parseByte2,
	//Half:      parseHalf,
	//Short:     parseShort,
	//Byte4:     parseByte4,
	//Word:      parseWord,
	//Long:      parseLong,
	//Zero:      parseZero,
	//VariantCC: parseVariantCC,
	//Attribute: parseAttribute,
}

func (st *Stmt) parseDirective(val string) error {
	d := Directive{
		name:    val,
		argTyps: directiveSet[val],
		src:     st.src,
		idx:     st.idx,
	}

	// 必要な引数の分だけコードを読み進めながらパース
	argTypIdx := 0
	d.skipUntilNextVal()

	for !d.isEOF() && argTypIdx < len(d.argTyps) {
		val, typ := d.nextVal()
		if d.argTyps[argTypIdx]&typ == 0 {
			return errors.New(fmt.Sprintf(ErrMsg, val[0]))
		}
		d.args = append(d.args, val)
		d.skipUntilNextVal()
		argTypIdx++
	}

	// その行に文字列が残っていたらエラー
	if argTypIdx != len(d.argTyps) {
		return errors.New("missing argument.")
	} else if !d.isEOF() {
		return errors.New(fmt.Sprintf(ErrMsg, d.src[d.idx]))
	}

	if d.isSection() {
		st.section = d.name
	}
	st.dir = &d
	return nil
}
