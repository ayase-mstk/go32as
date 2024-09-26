package elf32

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"

	"github.com/ayase-mstk/go32as/src/parse"
)

// RV32I命令に対応する構造体
type Instruction struct {
	opcode int
	funct3 int
	funct7 int
}

// RV32I命令セットの全命令を定義するマップ
var instructionMap = map[string]Instruction{
	// R-type instructions
	"add":  {opcode: 0b0110011, funct3: 0b000, funct7: 0b0000000},
	"sub":  {opcode: 0b0110011, funct3: 0b000, funct7: 0b0100000},
	"xor":  {opcode: 0b0110011, funct3: 0b100, funct7: 0b0000000},
	"or":   {opcode: 0b0110011, funct3: 0b110, funct7: 0b0000000},
	"and":  {opcode: 0b0110011, funct3: 0b111, funct7: 0b0000000},
	"sll":  {opcode: 0b0110011, funct3: 0b001, funct7: 0b0000000},
	"srl":  {opcode: 0b0110011, funct3: 0b101, funct7: 0b0000000},
	"sra":  {opcode: 0b0110011, funct3: 0b101, funct7: 0b0100000},
	"slt":  {opcode: 0b0110011, funct3: 0b010, funct7: 0b0000000},
	"sltu": {opcode: 0b0110011, funct3: 0b011, funct7: 0b0000000},

	// I-type instructions
	"addi":  {opcode: 0b0010011, funct3: 0b000, funct7: 0}, // funct7は不要
	"xori":  {opcode: 0b0010011, funct3: 0b100, funct7: 0}, // funct7は不要
	"ori":   {opcode: 0b0010011, funct3: 0b110, funct7: 0}, // funct7は不要
	"andi":  {opcode: 0b0010011, funct3: 0b111, funct7: 0}, // funct7は不要
	"slli":  {opcode: 0b0010011, funct3: 0b001, funct7: 0b0000000},
	"srli":  {opcode: 0b0010011, funct3: 0b101, funct7: 0b0000000},
	"srai":  {opcode: 0b0010011, funct3: 0b101, funct7: 0b0100000},
	"slti":  {opcode: 0b0010011, funct3: 0b010, funct7: 0}, // funct7は不要
	"sltiu": {opcode: 0b0010011, funct3: 0b011, funct7: 0}, // funct7は不要
	"lb":    {opcode: 0b0000011, funct3: 0b000, funct7: 0}, // funct7は不要
	"lh":    {opcode: 0b0000011, funct3: 0b001, funct7: 0}, // funct7は不要
	"lw":    {opcode: 0b0000011, funct3: 0b010, funct7: 0}, // funct7は不要
	"lbu":   {opcode: 0b0000011, funct3: 0b100, funct7: 0}, // funct7は不要
	"lhu":   {opcode: 0b0000011, funct3: 0b101, funct7: 0}, // funct7は不要

	// S-type instructions
	"sb": {opcode: 0b0100011, funct3: 0b000, funct7: 0}, // funct7は不要
	"sh": {opcode: 0b0100011, funct3: 0b001, funct7: 0}, // funct7は不要
	"sw": {opcode: 0b0100011, funct3: 0b010, funct7: 0}, // funct7は不要

	// B-type instructions
	"beq":  {opcode: 0b1100011, funct3: 0b000, funct7: 0}, // funct7は不要
	"bne":  {opcode: 0b1100011, funct3: 0b001, funct7: 0}, // funct7は不要
	"blt":  {opcode: 0b1100011, funct3: 0b100, funct7: 0}, // funct7は不要
	"bge":  {opcode: 0b1100011, funct3: 0b101, funct7: 0}, // funct7は不要
	"bltu": {opcode: 0b1100011, funct3: 0b110, funct7: 0}, // funct7は不要
	"bgeu": {opcode: 0b1100011, funct3: 0b111, funct7: 0}, // funct7は不要

	// U-type instructions
	"lui":   {opcode: 0b0110111, funct3: 0, funct7: 0}, // funct3、funct7は不要
	"auipc": {opcode: 0b0010111, funct3: 0, funct7: 0}, // funct3、funct7は不要

	// J-type instructions
	"jal":  {opcode: 0b1101111, funct3: 0, funct7: 0},     // funct3、funct7は不要
	"jalr": {opcode: 0b1100111, funct3: 0b000, funct7: 0}, // funct7は不要

	"ecall":  {opcode: 0b1110011, funct3: 0, funct7: 0}, // 特殊命令
	"ebreak": {opcode: 0b1110011, funct3: 0, funct7: 0}, // 特殊命令
}

var RegisterEncode = map[string]int{
	"x0": 0, "zero": 0,
	"x1": 1, "ra": 1,
	"x2": 2, "sp": 2,
	"x3": 3, "gp": 3,
	"x4": 4, "tp": 4,
	"x5": 5, "t0": 5,
	"x6": 6, "t1": 6,
	"x7": 7, "t2": 7,
	"x8": 8, "s0": 8, "fp": 8,
	"x9": 9, "s1": 9,
	"x10": 10, "a0": 10,
	"x11": 11, "a1": 11,
	"x12": 12, "a2": 12,
	"x13": 13, "a3": 13,
	"x14": 14, "a4": 14,
	"x15": 15, "a5": 15,
	"x16": 16, "a6": 16,
	"x17": 17, "a7": 17,
	"x18": 18, "s2": 18,
	"x19": 19, "s3": 19,
	"x20": 20, "s4": 20,
	"x21": 21, "s5": 21,
	"x22": 22, "s6": 22,
	"x23": 23, "s7": 23,
	"x24": 24, "s8": 24,
	"x25": 25, "s9": 25,
	"x26": 26, "s10": 26,
	"x27": 27, "s11": 27,
	"x28": 28, "t3": 28,
	"x29": 29, "t4": 29,
	"x30": 30, "t5": 30,
	"x31": 31, "t6": 31,
}

// R型命令のエンコード
func encodeRType(instName string, rd, rs1, rs2 int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	return uint32(inst.funct7)<<25 |
		uint32(rs2)<<20 |
		uint32(rs1)<<15 |
		uint32(inst.funct3)<<12 |
		uint32(rd)<<7 |
		uint32(inst.opcode)
}

// I型命令のエンコード
func encodeIType(instName string, rd, rs1, imm int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	return uint32(imm&0xFFF)<<20 |
		uint32(rs1)<<15 |
		uint32(inst.funct3)<<12 |
		uint32(rd)<<7 |
		uint32(inst.opcode)
}

// S型命令のエンコード
func encodeSType(instName string, rs1, rs2, imm int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	imm11_5 := (imm >> 5) & 0x7F
	imm4_0 := imm & 0x1F
	return uint32(imm11_5)<<25 |
		uint32(rs2)<<20 |
		uint32(rs1)<<15 |
		uint32(inst.funct3)<<12 |
		uint32(imm4_0)<<7 |
		uint32(inst.opcode)
}

func encodeBType(instName string, rs1, rs2, imm int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	imm11 := (imm >> 11) & 0x1       // 即値の11ビット目
	imm4_1 := (imm >> 1) & 0xF       // 即値の4:1ビット目
	imm10_5 := (imm >> 5) & 0x3F     // 即値の10:5ビット目
	imm12 := (imm >> 12) & 0x1       // 即値の12ビット目

	return uint32(imm12)<<31 |
		uint32(imm10_5)<<25 |
		uint32(rs2)<<20 |
		uint32(rs1)<<15 |
		uint32(inst.funct3)<<12 |
		uint32(imm4_1)<<8 |
		uint32(imm11)<<7 |
		uint32(inst.opcode)
}

// U型命令のエンコーディング関数
func encodeUType(instName string, rd, imm int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	return uint32(imm)<<12 |
		uint32(rd)<<7 |
		uint32(inst.opcode)
}

// J型命令のエンコーディング関数
func encodeJType(instName string, rd, imm int) uint32 {
	inst := instructionMap[instName] // 命令名に基づいてインストラクション情報を取得
	imm20 := (imm >> 20) & 0x1       // 即値の20ビット目
	imm10_1 := (imm >> 1) & 0x3FF    // 即値の10:1ビット目
	imm11 := (imm >> 11) & 0x1       // 即値の11ビット目
	imm19_12 := (imm >> 12) & 0xFF   // 即値の19:12ビット目

	return uint32(imm20)<<31 |
		uint32(imm19_12)<<12 |
		uint32(imm11)<<20 |
		uint32(imm10_1)<<21 |
		uint32(rd)<<7 |
		uint32(inst.opcode)
}

func (e *Elf32) resolveImm(val string) int {
	// 即値の場合そのまま返す
	if parse.IsImmediate(val) {
		n, _ := strconv.Atoi(val)
		return n
	}

	// symbolの場合
	idx := e.symtbl.idx[val]
	sym := e.symtbl.symtbls[idx]
	return int(sym.value)
}

func changeLoadInstruction(opecode string, operands *[]string) {
	switch opecode {
	case "lb", "lh", "lw", "lbu", "lhu":
		(*operands)[1], (*operands)[2] = (*operands)[2], (*operands)[1]
	default:
		return
	}
}

func dataEncode(file *os.File, stmts []parse.Stmt) {
	for _, stmt := range stmts {
		switch stmt.Dir().Name() {
		case ".string", ".asciz":
			data := strings.Trim(stmt.Dir().Args()[0], "\"")
			file.Write([]byte(data))
		case ".byte":
			// overflowはパーサーで処理済みと仮定
			data, _ := strconv.Atoi(stmt.Dir().Args()[0])
			data8 := int8((data+(1<<7))%(1<<8) - (1 << 7))
			binary.Write(file, binary.LittleEndian, data8)
		case ".2byte", ".half", ".short":
			data, _ := strconv.Atoi(stmt.Dir().Args()[0])
			data16 := int16((data+(1<<15))%(1<<16) - (1 << 15))
			binary.Write(file, binary.LittleEndian, data16)
		case ".4byte", ".word":
			data, _ := strconv.Atoi(stmt.Dir().Args()[0])
			data32 := int32((data+(1<<31))%(1<<32) - (1 << 31))
			binary.Write(file, binary.LittleEndian, data32)
		}
	}
}

func encodeSymtblEntries(file *os.File, symtbls []Elf32SymtblEntry) error {
	for _, entry := range symtbls {
		err := binary.Write(file, binary.LittleEndian, entry.name)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.LittleEndian, entry.value)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.LittleEndian, entry.size)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.LittleEndian, entry.info)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.LittleEndian, entry.other)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.LittleEndian, entry.shndx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Elf32) encodeAttributes(file *os.File) error {
	// .riscv.attributes section
	err := binary.Write(file, binary.LittleEndian, e.attr.FormatVersion)
	if err != nil {
		return err
	}

	for _, subsection := range e.attr.VendorSections {
		err = binary.Write(file, binary.LittleEndian, subsection.Length)
		if err != nil {
			return err
		}

		// ベンダー名は NTBS なので、文字列をバイト配列に変換し、null 終端を追加
		vendorNameBytes := append([]byte(subsection.VendorName), 0) // null 終端を追加
		err = binary.Write(file, binary.LittleEndian, vendorNameBytes)
		if err != nil {
			return err
		}

		for _, subsubsection := range subsection.SubSubSections {
			err = binary.Write(file, binary.LittleEndian, subsubsection.Tag)
			if err != nil {
				return err
			}
			err = binary.Write(file, binary.LittleEndian, subsubsection.Length)
			if err != nil {
				return err
			}

			for _, attr := range subsubsection.Attributes {
				err = binary.Write(file, binary.LittleEndian, attr.Tag)
				if err != nil {
					return err
				}

				// attr.Value は ULEB128 または NTBS なので、タイプに応じて処理
				switch v := attr.Value.(type) {
				case ULEB128:
					// ULEB128 のエンコードを実行する
					encodedValue := encodeULEB128(v)
					_, err = file.Write(encodedValue)
					if err != nil {
						return err
					}
				case string:
					// NTBS のエンコードを実行
					valueBytes := append([]byte(v), 0) // null 終端を追加
					_, err = file.Write(valueBytes)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func encodeULEB128(value ULEB128) []byte {
	var buffer bytes.Buffer

	for {
		// 最上位ビットを0にして残りのビットを取る
		encodedByte := byte(value & 0x7F) // 7ビットを取得
		value >>= 7                       // 次の7ビットにシフト

		// 次のバイトがある場合、最上位ビットをセット
		if value != 0 {
			encodedByte |= 0x80 // 最上位ビットを1に設定
		}

		buffer.WriteByte(encodedByte) // エンコードされたバイトをバッファに追加

		// すべてのビットがエンコードされた場合は終了
		if value == 0 {
			break
		}
	}

	return buffer.Bytes() // エンコードされたバイトスライスを返す
}

func (e *Elf32) WriteToFile() error {
	file, err := os.Create("output.o")
	if err != nil {
		return err
	}
	defer file.Close()

	// ELF header
	err = binary.Write(file, binary.LittleEndian, e.ehdr)
	if err != nil {
		return err
	}

	// .text section
	for _, stmt := range e.sections.entry[".text"].stmts {
		var data uint32
		opcode := stmt.Op().Opecode()
		oprands := stmt.Op().Operands()
		switch stmt.Op().OpcType() {
		case parse.RType:
			rd := RegisterEncode[oprands[0]]
			rs1 := RegisterEncode[oprands[1]]
			rs2 := RegisterEncode[oprands[2]]
			data = encodeRType(opcode, rd, rs1, rs2)
		case parse.IType:
			if opcode == "ecall" || opcode == "ebreak" {
				data = encodeIType(opcode, 0, 0, 0)
				break
			}
			changeLoadInstruction(opcode, &oprands)
			rd := RegisterEncode[oprands[0]]
			rs1 := RegisterEncode[oprands[1]]
			imm := e.resolveImm(oprands[2])
			data = encodeIType(opcode, rd, rs1, imm)
		case parse.SType:
			rs1 := RegisterEncode[oprands[0]]
			imm := e.resolveImm(oprands[1])
			rs2 := RegisterEncode[oprands[2]]
			data = encodeSType(opcode, rs1, rs2, imm)
		case parse.BType:
			// 最適化があるようなので、そのまま計算するようなことはできなさそう。
			rs1 := RegisterEncode[oprands[0]]
			rs2 := RegisterEncode[oprands[1]]
			imm := e.resolveImm(oprands[2])
			data = encodeBType(opcode, rs1, rs2, imm)
		case parse.UType:
			rd, _ := RegisterEncode[oprands[0]]
			imm := e.resolveImm(oprands[1])
			data = encodeUType(opcode, rd, imm)
		case parse.JType:
			rd, _ := RegisterEncode[oprands[0]]
			imm := e.resolveImm(oprands[1])
			data = encodeJType(opcode, rd, imm)
		}
		err = binary.Write(file, binary.LittleEndian, data)
		if err != nil {
			return err
		}
	}

	// .data section
	dataEncode(file, e.sections.entry[".data"].stmts)
	// .bss section
	dataEncode(file, e.sections.entry[".bss"].stmts)
	// .rodata section
	dataEncode(file, e.sections.entry[".rodata"].stmts)

	// .riscv.attributes section
	// size = 0x4c 0x13ツールチェーンより少ない
	err = e.encodeAttributes(file)
	if err != nil {
		return err
	}

	// .symtab
	// Elf32SymtblEntryにエンコードしなくてよい要素も入っているのでそのままエンコードできない
	err = encodeSymtblEntries(file, e.symtbl.symtbls)
	if err != nil {
		return err
	}

	// .strtab
	_, err = file.Write(e.strtbl.data)
	if err != nil {
		return err
	}
	// .shstrtab
	_, err = file.Write(e.shstrtbl.data)
	if err != nil {
		return err
	}

	// .rela.text
	for _, entry := range e.rela.entry {
		err = binary.Write(file, binary.LittleEndian, entry)
		if err != nil {
			return err
		}
	}

	// section header table
	for _, entry := range e.shdr.shdrs {
		err = binary.Write(file, binary.LittleEndian, entry)
		if err != nil {
			return err
		}
	}
	return nil
}
