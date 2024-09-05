package parse

import (
  "os"
  "bufio"
  "fmt"
  "errors"
)

const ErrMsg string = "junk at end of line, first unrecognized character is `%c'"

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

var delim []rune = []rune{
  0x09, // HT
  0x0b, // VT
  0x0c, // FF
  0x20, // SP
  0x23, // #
  0x2c, // ,
}

func isDelim(r rune) bool {
  for _, d := range delim {
    if r == d {
      return true
    }
  }
  return false
}

type StmtType int
const (
  DIRECTIVE StmtType = iota
  OPERATION
  UNKNOWN
  // LABELの有無はlabelSymbolを見て判断する
)


type Stmt struct {
  typ         StmtType
  op          *Operation
  dir         *Directive
  labelSymbol string
  row         int
  src         []rune
  idx         int
}

func (s *Stmt) Type() StmtType { return s.typ }
func (s *Stmt) Op() *Operation { return s.op }
func (s *Stmt) Dir() *Directive { return s.dir }
func (s *Stmt) LSymbol() string { return s.labelSymbol }
func (s *Stmt) Row() int { return s.row }

func (s *Stmt) setType() {
  if s.Op() != nil {
    s.typ = OPERATION
  } else if s.Dir() != nil {
    s.typ = DIRECTIVE
  } else {
    s.typ = UNKNOWN
  }
}

func (s *Stmt) skipUntilNextToken() {
  for ; s.idx < len(s.src); s.idx++ {
    c := s.src[s.idx]
    if c != ' ' && c != '\t' {
      return
    } else if s.src[s.idx] == '#' { // comment以降は飛ばす
      s.idx = len(s.src)-1
      return
    }
  }
}

func (s *Stmt) isEOF() bool {
  return s.idx == len(s.src)
}

func (s *Stmt) getToken() Token {
  start := s.idx
  for ; s.idx < len(s.src); s.idx++ {
    if isDelim(s.src[s.idx]) {
      break
    }
  }
  // この関数に入った場合かならずtokenがある
  val := string(s.src[start:s.idx])
  newTK := newToken(val)
  return newTK
}

func ParseLine(input []rune, row int) (Stmt, error) {
  stmt := Stmt{
    op:     nil,
    dir:    nil,
    row:    row,
    src:    input,
    idx:    0,
  }

  stmt.skipUntilNextToken()
  if stmt.isEOF() {
    return stmt, nil
  }
  tk := stmt.getToken()

  if tk.Type() == TLabel{
    // set symbol
    stmt.labelSymbol = tk.Val()
    stmt.skipUntilNextToken()
    if stmt.isEOF() {
      stmt.setType()
      return stmt, nil
    }
    tk = stmt.getToken()
  }
  stmt.skipUntilNextToken()

  // Stmtタイプで処理を分ける
  switch tk.Type() {
  case TDirective:
    // それぞれのパーサーを呼ぶ
    err := stmt.parseDirective()
    stmt.setType()
    return stmt, err
  case TOpecode:
    // それぞれのパーサーを呼ぶ
    err := stmt.parseOperation(tk.Val())
    stmt.setType()
    return stmt, err
  default:
    return stmt, errors.New(fmt.Sprintf(ErrMsg, tk.Val()[0]))
  }
}

func ParseFile(filename string) ([]Stmt, error) {
  var stmts []Stmt

  // ファイルをオープンします。
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer file.Close() // 関数が終了する際にファイルをクローズします。

  // バッファードリーダーを作成します。
  scanner := bufio.NewScanner(file)
  row := 1

  // ファイルの各行を読み込みます。
  for scanner.Scan() {
    line := scanner.Text() // 現在の行を取得します。
    newStmt, err := ParseLine([]rune(line), row)
    if err != nil {
      return nil, fmt.Errorf("%s:%d: Error: %s\n", filename, row, err.Error())
    }
    stmts = append(stmts, newStmt) // 行を処理します。
    row++
  }

  // 読み込み中にエラーが発生した場合はエラーを返します。
  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return stmts, nil
}
