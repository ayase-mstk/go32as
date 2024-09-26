package elf32

import (
	"fmt"
	"strconv"

	"github.com/ayase-mstk/go32as/src/parse"
)

type Elf32 struct {
	ehdr     Elf32Ehdr
	sections Elf32Sections
	attr     Elf32Attributes
	symtbl   Symtbl
	strtbl   Elf32Strtbl
	shstrtbl Elf32Shstrtbl
	rela     Rela
	shdr     Shdr
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
	elf.initAttributes()
	// section header を初期化する
	elf.initSectionHeader()
	// symbol table のindex0にからシンボルを追加
	elf.initSymbolTables()

	// 1周目
	for _, stmt := range stmts {
		var off Elf32Addr

		if stmt.LSymbol() != "" {
			if !elf.symtbl.exist(stmt.LSymbol()) {
				// まだシンボルテーブルになければ追加
				labelName := stmt.LSymbol()
				newSym := newSymbol(elf.strtbl.resolveIndex(labelName), elf.sections.resolveOffset(stmt.Section()), 0, createSymInfo(STB_LOCAL, STT_NOTYPE), elf.shdr.resolveShndx(stmt.Section()), stmt.Section())
				elf.symtbl.addSymbol(newSym, labelName)
			} else {
				// 既にシンボルテーブルに存在するラベル名だった場合
				// 他のセクションに同名のシンボルがあったらエラー
				if elf.symtbl.duplicateLabel(stmt.LSymbol(), stmt.Section(), elf.strtbl) {
					return elf, fmt.Errorf("%d: Error: symbol `%q' is already defined\n", stmt.Row(), stmt.LSymbol())
				}
				// 重複していなければ、まだセクションに属していない可能性があるので、設定する
				elf.symtbl.setSection(stmt.LSymbol(), stmt.Section())
			}
		}

		if stmt.Dir() != nil {
			elf.handleDirective(stmt)
			off = calcSize(stmt)
		} else if stmt.Op() != nil {
			// codeがtextセクション以外にあったらエラー
			if stmt.Section() != parse.Text {
				return elf, fmt.Errorf("%d: Error: unknown pseudo-op:%s\n", stmt.Row(), stmt.Op().Opecode())
			}
			off = 4
			section := elf.sections.entry[".text"]
			section.stmts = append(section.stmts, stmt)
			elf.sections.entry[".text"] = section
		}
		elf.sections.advanceOffset(stmt.Section(), off)
	}

	// 2周目
	// 外部シンボル解決
	elf.resolveOperationSymbol()
	elf.resolveSymbolShndx()
	elf.ResolveSectionRayout() // section header table 作成
	elf.resolveELFHeader()
	return elf, nil
}

func (e *Elf32) handleDirective(s parse.Stmt) {
	switch s.Dir().Name() {
	case ".section", ".text", ".data", ".rodata", ".bss":
		// section symbolはnameを持たない
		newSym := newSymbol(0, 0, 0, createSymInfo(STB_LOCAL, STT_SECTION), e.shdr.resolveShndx(s.Section()), s.Section())
		// もしすでに存在していれば追加されない
		e.symtbl.addSymbol(newSym, s.Section())

		break

	case ".align":
		e.shdr.setAddrAlign(s.Section(), s.Dir().Args()[0])
		break

	case ".file":
		newSym := newSymbol(e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_FILE), SHN_ABS, "")
		e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
		break

	case ".local":
		if e.symtbl.exist(s.Dir().Args()[0]) {
			e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(STB_LOCAL, (e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info&0x0F)))
		} else {
			newSym := newSymbol(e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, STT_NOTYPE), SHN_UNDEF, "")
			e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
		}
		break

	case ".globl":
		if e.symtbl.exist(s.Dir().Args()[0]) {
			e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(STB_LOCAL, (e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info&0x0F)))
		} else {
			newSym := newSymbol(e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_GLOBAL, STT_NOTYPE), SHN_UNDEF, "")
			e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
		}
		break

	case ".comm", ".common":
		// shndxをSHN_COMMONに設定する
		break

	case ".ident":
		break

	case ".size":
		break

	case ".string", ".asciz": // alias for string
		section := s.Section()
		if section == ".data" {
			section := e.sections.entry[".data"]
			section.stmts = append(section.stmts, s)
			e.sections.entry[".data"] = section
		} else if section == ".rodata" {
			section := e.sections.entry[".rodata"]
			section.stmts = append(section.stmts, s)
			e.sections.entry[".rodata"] = section
		}
		break

	case ".equ":
		val, _ := strconv.Atoi(s.Dir().Args()[1])
		if e.symtbl.exist(s.Dir().Args()[0]) {
			e.symtbl.setValue(s.Dir().Args()[0], Elf32Addr(val))
		} else {
			newSym := newSymbol(e.strtbl.resolveIndex(s.Dir().Args()[0]), Elf32Addr(val), 0, createSymInfo(STB_LOCAL, STT_NOTYPE), SHN_ABS, "")
			e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
		}
		break

	case ".type":
		// TODO: 現状.typeに数値しか設定できない
		symbolType, _ := symbolInfoTypes[s.Dir().Args()[1]]
		if e.symtbl.exist(s.Dir().Args()[0]) {
			e.symtbl.setInfo(s.Dir().Args()[0], createSymInfo(e.symtbl.symtbls[e.symtbl.idx[s.Dir().Args()[0]]].info>>4, symbolType))
		} else {
			newSym := newSymbol(e.strtbl.resolveIndex(s.Dir().Args()[0]), 0, 0, createSymInfo(STB_LOCAL, symbolType), SHN_UNDEF, "")
			e.symtbl.addSymbol(newSym, s.Dir().Args()[0])
		}
		break

	case ".byte", ".2byte", ".half", ".short", ".4byte", ".word":
		if s.Section() == ".data" {
			section := e.sections.entry[".data"]
			section.stmts = append(section.stmts, s)
			e.sections.entry[".data"] = section
		} else if s.Section() == ".bss" {
			section := e.sections.entry[".bss"]
			section.stmts = append(section.stmts, s)
			e.sections.entry[".bss"] = section
		} else if s.Section() == ".rodata" {
			section := e.sections.entry[".rodata"]
			section.stmts = append(section.stmts, s)
			e.sections.entry[".rodata"] = section
		}
		break
	}
}

func calcSize(s parse.Stmt) Elf32Addr {
	var off Elf32Addr = 0

	switch s.Dir().Name() {
	case ".string", ".asciz": // alias for string
		off = Elf32Addr(len(s.Dir().Args()[0]) - 2) // double quotationの分減らす
		break
	case ".byte":
		off = 1
		break
	case ".2byte", ".half", ".short":
		off = 2
		break
	case ".4byte", ".word":
		off = 4
		break
	default:
		break
	}
	return off
}

// テーブル処理一週目の後に実行
// 命令文中に出てくるシンボルを解決
func (e *Elf32) resolveOperationSymbol() {
	// .textセクションだけ見る
	entry, exists := e.sections.entry[".text"]
	if !exists {
		return
	}
	var off Elf32Addr = 0
	for _, stmt := range entry.stmts {
		// 命令文中にシンボル名が使用されて場合、それがローカルのシンボルテーブル中に存在するか確認
		symName := stmt.Op().RetIfSymbol()
		if len(symName) > 0 {
			// 存在しなければ外部シンボルなので外部シンボルとしてシンボルテーブルに追加する
			if !e.symtbl.exist(symName) {
				newSym := newSymbol(e.strtbl.resolveIndex(symName), 0, 0, createSymInfo(STB_GLOBAL, STT_NOTYPE), SHN_UNDEF, "")
				e.symtbl.addSymbol(newSym, symName)
			}
			// 命令文中にシンボルが使用されていれば、リロケーションエントリを作成する
			typ := resolveRelocType(*stmt.Op())
			e.rela.addRelaEntry(off, e.symtbl.idx[symName], typ, 0)
			e.rela.addRelaEntry(off, 0, RELAX, 0)
		}
		off += 4
	}
	if len(e.rela.entry) > 0 {
		relaTextSection := Elf32Shdr{
			ShName:      e.shstrtbl.resolveIndex(".rela.text"),
			ShType:      SHTRela,
			ShFlags:     SHFAlloc | SHFExecinstr,
			ShAddr:      0,
			ShOffset:    0,
			ShSize:      0,
			ShLink:      0,
			ShInfo:      0,
			ShAddralign: 4,
			ShEntsize:   12,
		}
		e.shdr.AddSection(relaTextSection, ".rela.text")
	}
}

func (e *Elf32) resolveSymbolShndx() {
	for i, sym := range e.symtbl.symtbls {
		// 空の場合は何もしない
		if sym.section == "" {
			continue
		}
		// shndxの値を設定し、テーブルに再代入する
		sym.shndx = e.shdr.resolveShndx(sym.section)
		e.symtbl.symtbls[i] = sym
	}
}

// ELFヘッダーの残りの変数を埋める
func (e *Elf32) resolveELFHeader() {
	e.ehdr.EShnum = Elf32Half(len(e.shdr.shdrs))            // sectionの数
	e.ehdr.EShstrndx = Elf32Half(e.shdr.shndx[".shstrtab"]) // section header内での.shstrtabのインデックス
}
