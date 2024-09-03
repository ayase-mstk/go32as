package parse

type DirectiveType int

type DirectiveToken struct {
  Token
  dtyp DirectiveType
}

func (dt *DirectiveToken) DType() DirectiveType { return dt.dtyp }

const (
  SECTION DirectiveType = iota
  SYMBOL
  MACRO
)

type Directive struct {
  Args  ArgFlag
  Type  DirectiveType
  // 引数自体がvalidかどうかはparseで判断
}

// directive
// map[string][]string
const (
  Align     = ".align"
  // P2Align   = ".p2align"
  // BAlign    = ".balign"
  File      = ".file"
  Globl     = ".globl"
  Local     = ".local"
  Comm      = ".comm"
  Common    = ".common"
  Ident     = ".ident"
  Section   = ".section"
  Size      = ".size"
  Text      = ".text"
  Data      = ".data"
  RoData    = ".rodata"
  Bss       = ".bss"
  String    = ".string"
  Asciz     = ".asciz"
  Equ       = ".equ"
  Macro     = ".macro"
  Endm      = ".endm"
  Type      = ".type"
  // Option     = ".option"
  Byte      = ".byte"
  Byte2     = ".2byte"
  Half      = ".half"
  Short     = ".short"
  Byte4     = ".4byte"
  Word      = ".word"
  Long      = ".long"
  // Float      = ".float"
  // DtprelWord = ".dtprelword"
  Zero      = ".zero"
  VariantCC = ".variant_cc"
  Attribute = ".attribute"
)

var directiveSet = map[string]struct{}{
  Align:      {},
  // P2Align:    {},
  // BAlign:     {},
  File:       {},
  Globl:      {},
  Local:      {},
  Comm:       {},
  Common:     {},
  Ident:      {},
  Section:    {},
  Size:       {},
  Text:       {},
  Data:       {},
  RoData:     {},
  Bss:        {},
  String:     {},
  Asciz:      {},
  Equ:        {},
  Macro:      {},
  Endm:       {},
  Type:       {},
  // Option:     {},
  Byte:       {},
  Byte2:      {},
  Half:       {},
  Short:      {},
  Byte4:      {},
  Word:       {},
  Long:       {},
  // Float:      {},
  // DtprelWord: {},
  Zero:       {},
  VariantCC:  {},
  Attribute:  {},
}

var directiveMap = map[string]func([]Token) bool{
  //Align:     parseAlign,
  //File:      parseFile,
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


