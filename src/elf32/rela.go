package elf32

import (
	"fmt"

	"github.com/ayase-mstk/go32as/src/parse"
)

type Rela struct {
	entry []Elf32Rela
	idx   map[string]int
}

type Elf32Rela struct {
	Off    Elf32Addr
	Info   Elf32Word // 再配置を行うシンボルテーブルのインデックスと、適用する再配置のタイプ
	Addend Elf32Sword
}

func (r *Rela) printRela() {
	for _, rela := range r.entry {
		fmt.Printf("off=%q\n", rela.Off)
		fmt.Printf("info.sym=%q\n", RelaSym(rela.Info))
		fmt.Printf("info.typ=%q\n", RelaType(rela.Info))
		fmt.Printf("addend=%q\n", rela.Addend)
	}
}

// Infoからそれぞれシンボルとタイプを抽出する関数
func RelaSym(i Elf32Word) Elf32Word {
	return i >> 8
}
func RelaType(i Elf32Word) RelocType {
	return RelocType(i)
}

// シンボルとタイプからInfoを作成する関数
func createRelaInfo(s int, t RelocType) Elf32Word {
	return Elf32Word(s<<8) + Elf32Word(t)
}

func (r *Rela) addRelaEntry(offset Elf32Addr, symbolIdx int, relocType RelocType, addend Elf32Sword) {
	info := createRelaInfo(symbolIdx, relocType)
	entry := Elf32Rela{
		Off:    offset,
		Info:   info,
		Addend: addend,
	}
	r.entry = append(r.entry, entry)
}

func resolveRelocType(op parse.Operation) RelocType {
	switch op.OpcType() {
	case parse.JType:
		return JAL

	case parse.BType:
		return BRANCH

	case parse.UType:
		if op.RelFunc() == "%hi" {
			return HI20
		} else if op.RelFunc() == "%pcrel_hi" {
			return PCREL_HI20
		}
		break

	case parse.IType:
		if op.RelFunc() == "%lo" {
			return LO12_I
		} else if op.RelFunc() == "%pcrel_lo" {
			return PCREL_LO12_I
		}
		break

	case parse.SType:
		if op.RelFunc() == "%lo" {
			return LO12_S
		} else if op.RelFunc() == "%pcrel_lo" {
			return PCREL_LO12_S
		}
		break
	default:
		break
	}
	return NONE
}

type RelocType int

const (
	// Enum values for different relocation types.
	NONE RelocType = iota // 0: No relocation

	// Both static and dynamic relocations
	R32          // 1: 32-bit relocation, static and dynamic
	R64          // 2: 64-bit relocation, static and dynamic
	RELATIVE     // 3: Adjust link address (B + A), dynamic
	COPY         // 4: Must be in executable; not allowed in shared libraries, dynamic
	JUMP_SLOT    // 5: PLT entry relocation, dynamic
	TLS_DTPMOD32 // 6: TLS module ID, 32-bit, dynamic
	TLS_DTPMOD64 // 7: TLS module ID, 64-bit, dynamic
	TLS_DTPREL32 // 8: TLS offset, 32-bit, dynamic
	TLS_DTPREL64 // 9: TLS offset, 64-bit, dynamic
	TLS_TPREL32  // 10: TLS thread pointer offset, 32-bit, dynamic
	TLS_TPREL64  // 11: TLS thread pointer offset, 64-bit, dynamic
	TLSDESC      // 12: TLS Descriptor, dynamic

	// Static relocations
	BRANCH       // 16: 12-bit PC-relative branch offset
	JAL          // 17: 20-bit PC-relative jump offset
	CALL         // 18: Deprecated, use CALL_PLT for 32-bit PC-relative function call
	CALL_PLT     // 19: 32-bit PC-relative function call (PIC)
	GOT_HI20     // 20: High 20 bits of 32-bit PC-relative GOT access
	TLS_GOT_HI20 // 21: High 20 bits of 32-bit PC-relative TLS IE GOT access
	TLS_GD_HI20  // 22: High 20 bits of PC-relative TLS GD GOT reference
	PCREL_HI20   // 23: High 20 bits of PC-relative reference
	PCREL_LO12_I // 24: Low 12 bits of PC-relative reference (I-type)
	PCREL_LO12_S // 25: Low 12 bits of PC-relative reference (S-type)
	HI20         // 26: High 20 bits of 32-bit absolute address (U-Type)
	LO12_I       // 27: Low 12 bits of absolute address (I-type)
	LO12_S       // 28: Low 12 bits of absolute address (S-type)
	TPREL_HI20   // 29: High 20 bits of TLS thread pointer offset
	TPREL_LO12_I // 30: Low 12 bits of TLS thread pointer offset (I-type)
	TPREL_LO12_S // 31: Low 12 bits of TLS thread pointer offset (S-type)
	TPREL_ADD    // 32: TLS thread pointer addition

	// Other types (static)
	ADD8        // 33: 8-bit label addition
	ADD16       // 34: 16-bit label addition
	ADD32       // 35: 32-bit label addition
	ADD64       // 36: 64-bit label addition
	SUB8        // 37: 8-bit label subtraction
	SUB16       // 38: 16-bit label subtraction
	SUB32       // 39: 32-bit label subtraction
	SUB64       // 40: 64-bit label subtraction
	GOT32_PCREL // 41: PC-relative GOT entry reference, 32-bit
	ALIGN       // 43: Alignment statement
	RVC_BRANCH  // 44: 8-bit PC-relative branch offset
	RVC_JUMP    // 45: 11-bit PC-relative jump offset
	RELAX       // 51: Instruction relaxation
	SUB6        // 52: 6-bit local label subtraction
	SET6        // 53: 6-bit local label assignment
	SET8        // 54: 8-bit local label assignment
	SET16       // 55: 16-bit local label assignment
	SET32       // 56: 32-bit local label assignment
	IRELATIVE   // 58: Relocation against non-preemptible ifunc symbol
	PLT32       // 59: 32-bit relative offset to a function or its PLT entry
)
