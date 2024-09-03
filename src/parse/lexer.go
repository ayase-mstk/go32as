package parse

import (
  "os"
  "bufio"
  "strings"
  "fmt"
  "errors"
  "github.com/ayase-mstk/go32as/src/utils"
)

var delim []rune = []rune{
  0x09, // HT
  0x0a, // LF
  0x0b, // VT
  0x0c, // FF
  0x0d, // CR
  0x20, // SP
  0x23, // #
  // 0x2a, // *
  0x2c, // ,
  // 0x2f, // /
}

func isDelim(r rune) bool {
  for _, d := range delim {
    if r == d {
      return true
    }
  }
  return false
}

func isLabel(literal string) bool {
  // ラベルかどうか
  if !strings.HasSuffix(literal, ":") {
    return false
  }
  literal = literal[:len(literal)-1]
  // すべて数値
  if utils.IsNumericStr(literal[:len(literal)-1]) {
    return true;
  }
  // 接頭辞はalphabetかアンダーバー
  if utils.IsAlpha(literal[0]) ||  literal[0] == '_'  || literal[0] == '.'{
    literal = literal[1:]
    for i := 0; i < len(literal); i++ {
      if !(utils.IsAlpha(literal[i]) || utils.IsNumeric(literal[i]) || literal[i] == '_' || literal[i] == '.') {
        return false
      }
    }
  } else {
    return false
  }
  return true
}

func isDirective(val string) bool {
  _, exists := directiveSet[val]
  return exists
}

func isOpecode(val string) bool {
  _, exists := OpecodeSet[val]
  return exists
}

func whichToken(val string) TokenType {
  if isLabel(val) {
    return LABEL
  } else if isDirective(val) {
    return DIRECTIVE
  } else if isOpecode(val) {
    return OPECODE
  } else {
    return UNKNOWN
  }
}

type syntax struct {
  isLiteral bool
  hasLabel  bool
  hasDir    bool
  hasOpcode bool
}

func isValidSyntax(tk Token, st *syntax) error {
  // syntax error handle
  switch tk.Type() {
  case LABEL:
    if st.hasLabel {
      return errors.New("Multiple labels found on the same line. Only one label is allowed per line.")
    }
    if st.hasDir || st.hasOpcode {
      return fmt.Errorf("junk at end of line, first unrecognized character is `%c'", tk.Val()[0])
    }
    st.hasLabel = true
    break;

  case DIRECTIVE:
    if st.hasDir || st.hasOpcode {
      return fmt.Errorf("junk at end of line, first unrecognized character is `%c'", tk.Val()[0])
    }
    st.hasDir = true
    break;

  case OPECODE:
    if st.hasDir || st.hasOpcode {
      return fmt.Errorf("junk at end of line, first unrecognized character is `%c'", tk.Val()[0])
    }
    st.hasOpcode = true
    break;

  default:
    break;
  }
  return nil
}

func addNewToken(start, end, row int, input []rune, tokens *[]IToken, st *syntax) error {
  if end - start > 1 {
    val := string(input[start:end])
    newTK := newToken(whichToken(val), val, row)
    *tokens = append(*tokens, newTK)
    err := isValidSyntax(newTK, st)
    if err != nil {
      return err
    }
  }
  return nil
}

func LexerLine(input []rune, row int) ([]IToken, error) {
  var tokens []IToken
  start      := 0
  st := syntax{
    isLiteral: false,
    hasLabel:  false,
    hasDir:    false,
    hasOpcode: false,
  }

  for i := 0; i < len(input); i++ {
    // literal
    if '"' == input[i] && !st.isLiteral {
      st.isLiteral = true
    } else if '"' == input[i] && st.isLiteral {
      st.isLiteral = false
    } else if st.isLiteral {
      continue
    }

    if isDelim(input[i]) {
      err := addNewToken(start, i, row, input, &tokens, &st)
      if err != nil {
        return nil, err
      }
      // commentはそれ以降読み飛ばす
      if '#' == input[i] {
        start = len(input)-1
        break
      }
      start = i+1
    }
  }
  err := addNewToken(start, len(input), row, input, &tokens, &st)
  if err != nil {
    return nil, err
  }
  // printTokens(tokens)
  return tokens, nil
}

func LexerFile(filename string) ([]IToken, error) {
  var tokens []IToken

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
    newTokens, err := LexerLine([]rune(line), row)
    if err != nil {
      return nil, fmt.Errorf("%s:%d: Error: %s\n", filename, row, err.Error())
    }
    tokens = append(tokens, newTokens...)     // 行を処理します。
    row++
  }

  // 読み込み中にエラーが発生した場合はエラーを返します。
  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return tokens, nil
}
