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
  // idをつける必要あり。
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

type Shdr struct {
  shdrs []Elf32Shdr
  shnx   map[string]int
}

func (s *Shdr) resolveShnx(name string) Elf32Half {
  return Elf32Half(s.shnx[name])
}

func (s *Shdr) printSectionHeader(strtbl Elf32Strtbl) {
  for _, sh := range s.shdrs {
    end := int(sh.ShName)
    for end < len(strtbl.data) && strtbl.data[end] != 0 {
      end++
    }
    fmt.Printf("Name=%q\n", strtbl.data[sh.ShName:end])
    fmt.Println("Type=", sh.ShType)
    fmt.Println("Flags=", sh.ShFlags)
    fmt.Println("Addr=", sh.ShAddr)
    fmt.Println("Offset=", sh.ShOffset)
    fmt.Println("Size=", sh.ShSize)
    fmt.Println("Link=", sh.ShLink)
    fmt.Println("Info=", sh.ShInfo)
    fmt.Println("Addralign=", sh.ShAddralign)
    fmt.Println("EntSize=", sh.ShEntsize)
    fmt.Println("")
  }
}

func (s *Shdr) AddSection(shdr Elf32Shdr, name string) {
    // shdrs に新しいセクションを追加
    s.shdrs = append(s.shdrs, shdr)

    // shnx にセクション名をキー、shdrs のインデックスを値として追加
    s.shnx[name] = len(s.shdrs) - 1
}

func (e *Elf32) initSectionHeader() {
  e.shdr.shnx = make(map[string]int)

  nullSection := Elf32Shdr{
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
  e.shdr.AddSection(nullSection, "")

  textSection := Elf32Shdr{
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
  e.shdr.AddSection(textSection, ".text")

  dataSection := Elf32Shdr{
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
  e.shdr.AddSection(dataSection, ".data")

  bssSection := Elf32Shdr{
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
  e.shdr.AddSection(bssSection, ".bss")

  rodataSection := Elf32Shdr{
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
  e.shdr.AddSection(rodataSection, ".rodata")
}

func (s *Shdr) setSectionAlignment(name, alignStr string) {
  idx := s.shnx[name]
  align, _ := strconv.Atoi(alignStr)
  s.shdrs[idx].ShAddralign = Elf32Word(math.Pow(2, float64(align)))
}
