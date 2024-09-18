package elf32

import (
	"fmt"
	"math"
	"strconv"
)

// ELFセクションタイプ
const (
	SHTNull            = 0  // 無効なセクション
	SHTProgbits        = 1  // プログラム情報
	SHTSymtab          = 2  // シンボルテーブル
	SHTStrtab          = 3  // 文字列テーブル
	SHTRela            = 4  // 明示的な追加情報付き再配置エントリ
	SHTHash            = 5  // ハッシュテーブル
	SHTDynamic         = 6  // 動的リンク情報
	SHTNote            = 7  // その他の情報
	SHTNobits          = 8  // ファイルに含まれないがメモリ上に必要な情報
	SHTRel             = 9  // 再配置エントリ
	SHTShlib           = 10 // 保留（意味は定義されていない）
	SHTDynsym          = 11 // 動的シンボルテーブル
	SHTRiscvAttributes = 0x70000003
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

type Shdr struct {
	shdrs []Elf32Shdr
	shndx map[string]int
}

// なければゼロが返る
func (s *Shdr) resolveShndx(name string) Elf32Half {
	return Elf32Half(s.shndx[name])
}

func (s *Shdr) printSectionHeader(shstrtbl Elf32Shstrtbl) {
	for _, sh := range s.shdrs {
		end := int(sh.ShName)
		for end < len(shstrtbl.data) && shstrtbl.data[end] != 0 {
			end++
		}
		fmt.Printf("Name=%q\n", shstrtbl.data[sh.ShName:end])
		fmt.Printf("Type=%#x\n", sh.ShType)
		fmt.Printf("Flags=%#x\n", sh.ShFlags)
		fmt.Printf("Addr=%#x\n", sh.ShAddr)
		fmt.Printf("Offset=%#x\n", sh.ShOffset)
		fmt.Printf("Size=%#x\n", sh.ShSize)
		fmt.Printf("Link=%#x\n", sh.ShLink)
		fmt.Printf("Info=%#x\n", sh.ShInfo)
		fmt.Printf("Addralign=%#x\n", sh.ShAddralign)
		fmt.Printf("EntSize=%#x\n", sh.ShEntsize)
		fmt.Println("")
	}
}

func (s *Shdr) AddSection(shdr Elf32Shdr, name string) {
	// shdrs に新しいセクションを追加
	s.shdrs = append(s.shdrs, shdr)

	// shndx にセクション名をキー、shdrs のインデックスを値として追加
	s.shndx[name] = len(s.shdrs) - 1
}

func (e *Elf32) initSectionHeader() {
	e.shdr.shndx = make(map[string]int)

	nullSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(""),
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
		ShName:      e.shstrtbl.resolveIndex(".text"),
		ShType:      SHTProgbits,
		ShFlags:     SHFAlloc | SHFExecinstr,
		ShAddr:      0,
		ShOffset:    0x34,
		ShSize:      0,
		ShLink:      0,
		ShInfo:      0,
		ShAddralign: 4,
		ShEntsize:   0,
	}
	e.shdr.AddSection(textSection, ".text")

	dataSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(".data"),
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
		ShName:      e.shstrtbl.resolveIndex(".bss"),
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

	riscvSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(".riscv.attributes"),
		ShType:      SHTRiscvAttributes,
		ShFlags:     0,
		ShAddr:      0,
		ShOffset:    0,
		ShSize:      0x5f,
		ShLink:      0,
		ShInfo:      0,
		ShAddralign: 1,
		ShEntsize:   0,
	}
	e.shdr.AddSection(riscvSection, ".riscv.attributes")

	//rodataSection := Elf32Shdr{
	//    ShName:      e.shstrtbl.resolveIndex(".rodata"),
	//    ShType:      SHTProgbits,
	//    ShFlags:     SHFAlloc,
	//    ShAddr:      0,
	//    ShOffset:    0,
	//    ShSize:      0,
	//    ShLink:      0,
	//    ShInfo:      0,
	//    ShAddralign: 4,
	//    ShEntsize:   0,
	//}
	//e.shdr.AddSection(rodataSection, ".rodata")

	symSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(".symtab"),
		ShType:      SHTSymtab,
		ShFlags:     0,
		ShAddr:      0,
		ShOffset:    0,
		ShSize:      0,
		ShLink:      0,
		ShInfo:      0,
		ShAddralign: 4,
		ShEntsize:   0,
	}
	e.shdr.AddSection(symSection, ".symtab")

	strSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(".strtab"),
		ShType:      SHTStrtab,
		ShFlags:     0,
		ShAddr:      0,
		ShOffset:    0,
		ShSize:      0,
		ShLink:      0,
		ShInfo:      0,
		ShAddralign: 1,
		ShEntsize:   0,
	}
	e.shdr.AddSection(strSection, ".strtab")

	shstrSection := Elf32Shdr{
		ShName:      e.shstrtbl.resolveIndex(".shstrtab"),
		ShType:      SHTStrtab,
		ShFlags:     0,
		ShAddr:      0,
		ShOffset:    0,
		ShSize:      0,
		ShLink:      0,
		ShInfo:      0,
		ShAddralign: 1,
		ShEntsize:   0,
	}
	e.shdr.AddSection(shstrSection, ".shstrtab")
}

func (s *Shdr) setAddrAlign(name, alignStr string) {
	idx := s.shndx[name]
	align, _ := strconv.Atoi(alignStr)
	s.shdrs[idx].ShAddralign = Elf32Word(math.Pow(2, float64(align)))
}

func (s *Shdr) setSize(name string, size Elf32Word) {
	idx := s.shndx[name]
	s.shdrs[idx].ShSize = size
}

func (s *Shdr) setOffset(name string, off Elf32Off) {
	idx := s.shndx[name]
	s.shdrs[idx].ShOffset = off
}

func (s *Shdr) setLink(name string, link Elf32Word) {
	idx := s.shndx[name]
	s.shdrs[idx].ShLink = link
}

func (s *Shdr) setInfo(name string, info Elf32Word) {
	idx := s.shndx[name]
	s.shdrs[idx].ShInfo = info
}

func (s *Shdr) setEntsize(name string, es Elf32Word) {
	idx := s.shndx[name]
	s.shdrs[idx].ShEntsize = es
}

func (s *Shdr) getOffset(name string) Elf32Off {
	idx := s.shndx[name]
	return s.shdrs[idx].ShOffset
}

func (s *Shdr) getSize(name string) Elf32Word {
	idx := s.shndx[name]
	return s.shdrs[idx].ShSize
}

func (s *Shdr) getEntsize(name string) Elf32Word {
	idx := s.shndx[name]
	return s.shdrs[idx].ShEntsize
}
