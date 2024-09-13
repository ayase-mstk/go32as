package elf32

import (
    "fmt"
    "strconv"

    "github.com/ayase-mstk/go32as/src/parse"
)

type Elf32 struct {
  ehdr      Elf32Ehdr
  symtbl    Symtbl
  strtbl    Elf32Strtbl
  shstrtbl  Elf32Shstrtbl
  shdr      Shdr
  sections  Elf32Sections
}

type Elf32Sections struct {
  text      []parse.Stmt
  data      []parse.Stmt
  rodata    []parse.Stmt
  bss       []parse.Stmt
}

func (e *Elf32) PrintAll() {
  fmt.Println("===================================")
  fmt.Println("============ ELF Header ===========")
  fmt.Println("===================================")
  e.ehdr.printELFHeader()
  fmt.Println("")
  fmt.Println("===================================")
  fmt.Println("=========== Symbol Table ==========")
  fmt.Println("===================================")
  e.symtbl.printSymbolTable(e.strtbl)
  fmt.Println("")
  fmt.Println("===================================")
  fmt.Println("========== Section Header =========")
  fmt.Println("===================================")
  e.shdr.printSectionHeader(e.shstrtbl)
}

/*
セクションヘッダーテーブルの初期化と、シンボルテーブルへのラベルとセクションの追加を行い、データ行とコード行を各セクションに分ける
*/
func PrepareElf32Tables(stmts []parse.Stmt) (Elf32, error) {
  var elf Elf32

  elf.initHeader()
  // section header を初期化する
  elf.initSectionHeader()
  // symbol table のindex0にからシンボルを追加
  elf.initSymbolTables()


  for _, stmt := range stmts {

    if stmt.LSymbol() != "" {
      if !elf.symtbl.SymbolExists(stmt.LSymbol()) {
        labelName := stmt.LSymbol()[:len(stmt.LSymbol())-1]
        newSym := newSymbol(stmt.LSymbol(), elf.strtbl.resolveIndex(labelName), 0, 0, createSymInfo(STB_LOCAL, STT_NOTYPE), elf.shdr.resolveShndx(stmt.Section()))
        elf.symtbl.addSymbol(newSym, stmt.Section())
      }
    }

    if stmt.Dir() != nil {
      elf.handleDirective(stmt)
    } else if stmt.Op() != nil {
      // codeがtextセクション以外にあったらエラー
      if stmt.Section() != parse.Text {
        return elf, fmt.Errorf("%d: Error: unknown pseudo-op:%s\n", stmt.Row(), stmt.Op().Opecode())
      }
      elf.sections.text = append(elf.sections.text, stmt)
    }
  }
  return elf, nil
}

func (e *Elf32) handleDirective(s parse.Stmt) {
  switch s.Dir().Name() {
  case ".section":
  case ".text":
  case ".data":
  case ".rodata":
  case ".bss":
    if !e.symtbl.SymbolExists(s.Section()) {
      newSym := newSymbol(s.Section(), e.shstrtbl.resolveIndex(s.Section()), 0, 0, createSymInfo(STB_LOCAL, STT_SECTION), e.shdr.resolveShndx(s.Section()))
      e.symtbl.addSymbol(newSym, s.Section())
    }

    break

  case ".align":
    e.shdr.setAddrAlign(s.Section(), s.Dir().Args()[0])
    break

  case ".file":
    newSym := newSymbol(s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_FILE), SHN_ABS)
    e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
    break

  case ".local":
    if e.symtbl.SymbolExists(s.Dir().Args()[0]) {
      e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(STB_LOCAL, (e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info & 0x0F)))
    } else {
      newSym := newSymbol(s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_NOTYPE), e.shdr.resolveShndx(s.Section()))
      e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
    }
    break

  case ".global":
    if e.symtbl.SymbolExists(s.Dir().Args()[0]) {
      e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(STB_LOCAL, (e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info & 0x0F)))
    } else {
      newSym := newSymbol(s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_GLOBAL, STT_NOTYPE), e.shdr.resolveShndx(s.Section()))
      e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
    }
    break

  case ".comm":
  case ".common":
    break

  case ".ident":
    break

  case ".size":
    break

  case ".string":
  case ".asciz": // alias for string
    section := s.Section()
    if section == ".data" {
      e.sections.data = append(e.sections.data, s)
    } else if section == ".rodata" {
      e.sections.rodata = append(e.sections.rodata, s)
    }
    break

  case ".equ":
    val, _ := strconv.Atoi(s.Dir().Args()[1])
    if e.symtbl.SymbolExists(s.Dir().Args()[0]) {
      e.symtbl.setValue(s.Dir().Args()[0], Elf32Addr(val))
    } else {
      newSym := newSymbol(s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), Elf32Addr(val), 0, createSymInfo(STB_LOCAL, STT_NOTYPE), e.shdr.resolveShndx(s.Section()))
      e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
    }
    break

  case ".type":
    symbolType := symbolInfoTypes[s.Dir().Args()[1]]
    if e.symtbl.SymbolExists(s.Dir().Args()[0]) {
      e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info >> 4, symbolType))
    } else {
      newSym := newSymbol(s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, symbolType), e.shdr.resolveShndx(s.Section()))
      e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
    }
    break

  case ".byte":
  case ".2byte":
  case ".half":
  case ".short":
  case ".4byte":
  case ".word":
    if s.Section() == ".data" {
      e.sections.data = append(e.sections.data, s)
    } else if s.Section() == ".bss" {
      e.sections.bss = append(e.sections.bss, s)
    } else if s.Section() == ".rodata" {
      e.sections.rodata = append(e.sections.rodata, s)
    }
    break
  }
}
