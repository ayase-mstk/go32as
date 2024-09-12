package elf32

import (
    "fmt"
    "strconv"

    "github.com/ayase-mstk/go32as/src/parse"
)

type Elf32 struct {
  ehdr      Elf32Ehdr
  symtbl    map[string]Elf32SymtblEntry
  strtbl    Elf32Strtbl
  shdr      map[string]Elf32Shdr
  sections   Elf32Sections
}

type Elf32Sections struct {
  text      []parse.Stmt
  data      []parse.Stmt
  rodata    []parse.Stmt
  bss       []parse.Stmt
}

func (e Elf32) PrintAll() {
  e.ehdr.printELFHeader()
  for key, _ := range e.symtbl {
    e.symtbl[key].printSymbolTable(e.strtbl)
  }
  for key, _ := range e.shdr {
    e.shdr[key].printSectionHeader(e.strtbl)
  }
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
      if _, exist := elf.symtbl[stmt.LSymbol()]; !exist {
        addSymbol(&elf.symtbl, stmt.LSymbol(), elf.strtbl.resolveIndex(stmt.LSymbol()), 0, 0, createSymInfo(STB_LOCAL, STT_NOTYPE), elf.resolveShnx(stmt.Section()))
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
    if _, exist := e.symtbl[s.Section()]; !exist {
      addSymbol(&e.symtbl, s.Section(), e.strtbl.resolveIndex(s.Section()), 0, 0, createSymInfo(STB_LOCAL, STT_SECTION), e.resolveShnx(s.Section()))
    }

    break

  case ".align":
    e.shdr[s.Section()] = setSectionAlignment(e.shdr[s.Section()], s.Dir().Args()[0])
    break

  case ".file":
    addSymbol(&e.symtbl, s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_FILE), SHN_ABS)
    break

  case ".local":
    if symtbl, exist := e.symtbl[s.Dir().Args()[0]]; exist {
      e.symtbl[s.Dir().Args()[0]] = setSymbolInfo(symtbl, createSymInfo(STB_LOCAL, (symtbl.info & 0x0F)))
    } else {
      addSymbol(&e.symtbl, s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_NOTYPE), e.resolveShnx(s.Section()))
    }
    break

  case ".global":
    if symtbl, exist := e.symtbl[s.Dir().Args()[0]]; exist {
      e.symtbl[s.Dir().Args()[0]] = setSymbolInfo(symtbl, createSymInfo(STB_GLOBAL, (symtbl.info & 0x0F)))
    } else {
      addSymbol(&e.symtbl, s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_GLOBAL, STT_NOTYPE), e.resolveShnx(s.Section()))
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
    if symtbl, exist := e.symtbl[s.Dir().Args()[0]]; exist {
      e.symtbl[s.Dir().Args()[0]] = setSymbolValue(symtbl, Elf32Addr(val))
    } else {
      addSymbol(&e.symtbl, s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), Elf32Addr(val), 0, createSymInfo(STB_LOCAL, STT_NOTYPE), e.resolveShnx(s.Section()))
    }
    break

  case ".type":
    symbolType := symbolInfoTypes[s.Dir().Args()[1]]
    if symtbl, exist := e.symtbl[s.Dir().Args()[0]]; exist {
      e.symtbl[s.Dir().Args()[0]] = setSymbolInfo(symtbl, createSymInfo(symtbl.info >> 4, symbolType))
    } else {
      addSymbol(&e.symtbl, s.Dir().Args()[0],  e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, symbolType), e.resolveShnx(s.Section()))
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
