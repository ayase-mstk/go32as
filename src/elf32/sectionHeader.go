package elf32

import (
  "fmt"
  "math"
  "strconv"
)

// ELFセクションタイプ
const (
	SHTNull     = 0  // 無効なセクション
	SHTProgbits = 1  // プログラム情報
	SHTSymtab   = 2  // シンボルテーブル
	SHTStrtab   = 3  // 文字列テーブル
	SHTRela     = 4  // 明示的な追加情報付き再配置エントリ
	SHTHash     = 5  // ハッシュテーブル
	SHTDynamic  = 6  // 動的リンク情報
	SHTNote     = 7  // その他の情報
	SHTNobits   = 8  // ファイルに含まれないがメモリ上に必要な情報
	SHTRel      = 9  // 再配置エントリ
	SHTShlib    = 10 // 保留（意味は定義されていない）
	SHTDynsym   = 11 // 動的シンボルテーブル
)

// セクションフラグ
const (
	SHFWrite     = 0x1 // セクションが書き込み可能
	SHFAlloc     = 0x2 // セクションがメモリにロードされる
	SHFExecinstr = 0x4 // セクションが実行可能な命令を含む
)

// ELF32セクションヘッダー構造体
type Elf32Shdr struct {
	ShName      Elf32Word // セクション名（文字列テーブルインデックス）
	ShType      Elf32Word // セクションのタイプ（例: SHTProgbits, SHTSymtab）
	ShFlags     Elf32Word // セクションフラグ（例: SHFWrite, SHFAlloc）
	ShAddr      Elf32Addr // セクションが配置されるメモリアドレス
	ShOffset    Elf32Off  // セクションのファイルオフセット
	ShSize      Elf32Word // セクションのサイズ（バイト単位）
	ShLink      Elf32Word // セクションリンク情報
	ShInfo      Elf32Word // 追加情報
	ShAddralign Elf32Word // セクションのメモリアライメント
	ShEntsize   Elf32Word // セクションのエントリサイズ（テーブルの場合）
}

func (e Elf32Shdr) printSectionHeader(strtbl Elf32Strtbl) {
  end := int(e.ShName)
  for end < len(strtbl.data) && strtbl.data[end] != 0 {
    end++
  }
  fmt.Printf("Name=%q\n", strtbl.data[e.ShName:end])
  fmt.Println("Type=", e.ShType)
  fmt.Println("Flags=", e.ShFlags)
  fmt.Println("Addr=", e.ShAddr)
  fmt.Println("Offset=", e.ShOffset)
  fmt.Println("Size=", e.ShSize)
  fmt.Println("Link=", e.ShLink)
  fmt.Println("Info=", e.ShInfo)
  fmt.Println("Addralign=", e.ShAddralign)
  fmt.Println("EntSize=", e.ShEntsize)
}

func (e *Elf32) initSectionHeader() {
  e.shdr = make(map[string]Elf32Shdr)
  e.shdr[""] = Elf32Shdr{
      ShName:      e.strtbl.resolveIndex(""),
      ShType:      SHTNull,
      ShFlags:     0,
      ShAddr:      0,
      ShOffset:    0,
      ShSize:      0,
      ShLink:      0,
      ShInfo:      0,
      ShAddralign: 0,
      ShEntsize:   0,
  }
  e.shdr[".text"] = Elf32Shdr{
      ShName:      e.strtbl.resolveIndex(".text"),
      ShType:      SHTProgbits,
      ShFlags:     SHFAlloc | SHFExecinstr,
      ShAddr:      0,
      ShOffset:    0,
      ShSize:      0,
      ShLink:      0,
      ShInfo:      0,
      ShAddralign: 4,
      ShEntsize:   0,
  }
  e.shdr[".data"] = Elf32Shdr{
      ShName:      e.strtbl.resolveIndex(".data"),
      ShType:      SHTProgbits,
      ShFlags:     SHFWrite | SHFAlloc,
      ShAddr:      0,
      ShOffset:    0,
      ShSize:      0,
      ShLink:      0,
      ShInfo:      0,
      ShAddralign: 4,
      ShEntsize:   0,
  }
  e.shdr[".bss"] = Elf32Shdr{
      ShName:      e.strtbl.resolveIndex(".bss"),
      ShType:      SHTNobits,
      ShFlags:     SHFWrite | SHFAlloc,
      ShAddr:      0,
      ShOffset:    0,
      ShSize:      0,
      ShLink:      0,
      ShInfo:      0,
      ShAddralign: 4,
      ShEntsize:   0,
  }
  e.shdr[".rodata"] = Elf32Shdr{
      ShName:      e.strtbl.resolveIndex(".rodata"),
      ShType:      SHTProgbits,
      ShFlags:     SHFAlloc,
      ShAddr:      0,
      ShOffset:    0,
      ShSize:      0,
      ShLink:      0,
      ShInfo:      0,
      ShAddralign: 4,
      ShEntsize:   0,
  }
}

func setSectionAlignment(s Elf32Shdr, alignStr string) Elf32Shdr {
  align, _ := strconv.Atoi(alignStr)
  s.ShAddralign = Elf32Word(math.Pow(2, float64(align)))
  return s
}

func (e *Elf32) resolveShnx(name string) Elf32Half {
  dataIdx := 0
  ret := 0
  for _, _ = range e.shdr {
    if string(e.strtbl.data[dataIdx:len(name)]) == name {
      break
    }
    ret++
    dataIdx++
  }
  return Elf32Half(ret)
}
