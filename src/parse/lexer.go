package parse

import (
  "os"
  "bufio"
  "strings"
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

func LexerLine(input []rune) []IToken {
  var tokens []IToken
  i         := 0
  start     := 0
  isLiteral := false

  for ; i < len(input); i++ {
    // literal
    if '"' == input[i] && !isLiteral {
      isLiteral = true
    } else if '"' == input[i] && isLiteral {
      isLiteral = false
    } else if isLiteral {
      continue
    }

    if isDelim(input[i]) {
      if i - start > 1 {
        val := string(input[start:i])
        tokens = append(tokens, newToken(whichToken(val), val))
      }
      // commentはそれ以降読み飛ばす
      if '#' == input[i] {
        start = len(input)-1
        break
      }
      start = i+1
    }
  }
  if len(input) - start > 1 {
    val := string(input[start:])
    tokens = append(tokens, newToken(whichToken(val), val))
  }

  // printTokens(tokens)
  return tokens
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

  // ファイルの各行を読み込みます。
  for scanner.Scan() {
    line := scanner.Text() // 現在の行を取得します。
    tokens = append(tokens, LexerLine([]rune(line))...)     // 行を処理します。
  }

  // 読み込み中にエラーが発生した場合はエラーを返します。
  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return tokens, nil
}
