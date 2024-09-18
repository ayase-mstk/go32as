package parse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	// R format
	ADD  = "add"
	SUB  = "sub"
	XOR  = "xor"
	OR   = "or"
	AND  = "and"
	SLL  = "sll"
	SRL  = "srl"
	SRA  = "sra"
	SLT  = "slt"
	SLTU = "sltu"

	// I format
	ADDI   = "addi"
	WORI   = "xori"
	ORI    = "ori"
	ANDI   = "andi"
	SLLI   = "slli"
	SRLI   = "srli"
	SRAI   = "srai"
	SLTI   = "slti"
	SLTIU  = "sltiu"
	LB     = "lb"
	LH     = "lh"
	LW     = "lw"
	LBU    = "lbu"
	LHU    = "lhu"
	JALR   = "jalr"
	ECALL  = "ecall"
	EBREAK = "ebreak"

	// S format
	SB = "sb"
	SH = "sh"
	SW = "sw"

	// B format
	BEQ  = "beq"
	BNE  = "bne"
	BLT  = "blt"
	BGE  = "bge"
	BLTU = "bltu"
	BGEU = "bgeu"

	// U format
	LUI   = "lui"
	AUIPC = "auipc"

	// J format
	JAL = "jal"
)

var RegisterSet = map[string]bool{
	"x0": true, "zero": true,
	"x1": true, "ra": true,
	"x2": true, "sp": true,
	"x3": true, "gp": true,
	"x4": true, "tp": true,
	"x5": true, "t0": true,
	"x6": true, "t1": true,
	"x7": true, "t2": true,
	"x8": true, "s0": true, "fp": true,
	"x9": true, "s1": true,
	"x10": true, "a0": true,
	"x11": true, "a1": true,
	"x12": true, "a2": true,
	"x13": true, "a3": true,
	"x14": true, "a4": true,
	"x15": true, "a5": true,
	"x16": true, "a6": true,
	"x17": true, "a7": true,
	"x18": true, "s2": true,
	"x19": true, "s3": true,
	"x20": true, "s4": true,
	"x21": true, "s5": true,
	"x22": true, "s6": true,
	"x23": true, "s7": true,
	"x24": true, "s8": true,
	"x25": true, "s9": true,
	"x26": true, "s10": true,
	"x27": true, "s11": true,
	"x28": true, "t3": true,
	"x29": true, "t4": true,
	"x30": true, "t5": true,
	"x31": true, "t6": true,
}

type OperandType int

const (
	REG OperandType = 1 << iota // 0x00000001
	IMM                         // 0x00000010
	LAB                         // 0x00000100
)

type OpecodeInfo struct {
	opcTyp  OpecodeType
	oprTyps []OperandType
	// 必要に応じて他のフィールドを追加（例: 関数ポインタ、コメントなど）
}

type OpecodeType int

const (
	RType OpecodeType = iota
	IType
	SType
	BType
	JType
	UType
)

var OpecodeMap = map[string]OpecodeInfo{
	ADD:  {RType, []OperandType{REG, REG, REG | IMM | LAB}},
	SUB:  {RType, []OperandType{REG, REG, REG}},
	XOR:  {RType, []OperandType{REG, REG, REG | IMM | LAB}},
	OR:   {RType, []OperandType{REG, REG, REG | IMM | LAB}},
	AND:  {RType, []OperandType{REG, REG, REG | IMM | LAB}},
	SLL:  {RType, []OperandType{REG, REG, REG | IMM}}, // なぜかこの3つの命令はラベルをとれない
	SRL:  {RType, []OperandType{REG, REG, REG | IMM}},
	SRA:  {RType, []OperandType{REG, REG, REG | IMM}},
	SLT:  {RType, []OperandType{REG, REG, REG | IMM | LAB}},
	SLTU: {RType, []OperandType{REG, REG, REG | IMM | LAB}},

	ADDI:  {IType, []OperandType{REG, REG, IMM | LAB}},
	WORI:  {IType, []OperandType{REG, REG, IMM | LAB}},
	ORI:   {IType, []OperandType{REG, REG, IMM | LAB}},
	ANDI:  {IType, []OperandType{REG, REG, IMM | LAB}},
	SLLI:  {IType, []OperandType{REG, REG, IMM}}, // なぜかこの3つの命令はラベルをとれない
	SRLI:  {IType, []OperandType{REG, REG, IMM}},
	SRAI:  {IType, []OperandType{REG, REG, IMM}},
	SLTI:  {IType, []OperandType{REG, REG, IMM | LAB}},
	SLTIU: {IType, []OperandType{REG, REG, IMM | LAB}},
	LB:    {IType, []OperandType{REG, IMM | LAB, REG}},
	LH:    {IType, []OperandType{REG, IMM | LAB, REG}},
	LW:    {IType, []OperandType{REG, IMM | LAB, REG}},
	LBU:   {IType, []OperandType{REG, IMM | LAB, REG}},
	LHU:   {IType, []OperandType{REG, IMM | LAB, REG}},
	JALR:  {IType, []OperandType{REG, REG, IMM | LAB}},
	//JALR:   {IType, []OperandType{REG, IMM | LAB, REG}},
	ECALL:  {IType, []OperandType{}},
	EBREAK: {IType, []OperandType{}},

	SB: {SType, []OperandType{REG, IMM | LAB, REG}},
	SH: {SType, []OperandType{REG, IMM | LAB, REG}},
	SW: {SType, []OperandType{REG, IMM | LAB, REG}},

	BEQ:  {BType, []OperandType{REG, REG, IMM | LAB}},
	BNE:  {BType, []OperandType{REG, REG, IMM | LAB}},
	BLT:  {BType, []OperandType{REG, REG, IMM | LAB}},
	BGE:  {BType, []OperandType{REG, REG, IMM | LAB}},
	BLTU: {BType, []OperandType{REG, REG, IMM | LAB}},
	BGEU: {BType, []OperandType{REG, REG, IMM | LAB}},

	LUI:   {UType, []OperandType{REG, IMM | LAB}},
	AUIPC: {UType, []OperandType{REG, IMM | LAB}},

	JAL: {JType, []OperandType{REG, IMM | LAB}},
}

type Operation struct {
	opcode   string
	info     OpecodeInfo
	operands []string
	relFunc  string
	src      []rune
	idx      int
}

func (o *Operation) Opecode() string        { return o.opcode }
func (o *Operation) Operands() []string     { return o.operands }
func (o *Operation) RelFunc() string        { return o.relFunc }
func (o *Operation) OpcType() OpecodeType   { return o.info.opcTyp }
func (o *Operation) OprType() []OperandType { return o.info.oprTyps }

// 命令文中にシンボルが出現していればそれを返す関数
func (o *Operation) RetIfSymbol() string {
	for i, typ := range o.info.oprTyps {
		if typ&LAB != 0 && !isRegister(o.operands[i]) && !isImmediate(o.operands[i]) {
			return o.operands[i]
		}
	}
	return ""
}
func (o Operation) printOperation() {
	for i := 0; i < len(o.Opecode()); i++ {
		fmt.Printf("opecodeint=%d\n", o.Opecode()[i])
	}
	fmt.Println("opecode=", o.Opecode())
	for i := 0; i < len(o.Operands()); i++ {
		fmt.Printf("operand[%d]=%s", i, o.Operands()[i])
	}
}

func (o *Operation) skipUntilNextOperand() {
	for ; o.idx < len(o.src); o.idx++ {
		c := o.src[o.idx]
		if o.src[o.idx] == '#' { // comment以降は飛ばす
			o.idx = len(o.src)
			return
		} else if c != ' ' && c != '\t' && c != ',' && c != '(' && c != ')' {
			return
		}
	}
}

func (o *Operation) isEOF() bool {
	return o.idx == len(o.src)
}

func (o *Operation) nextOperand() (string, OperandType) {
	isLiteral := false
	hasRelFunc := false
	start := o.idx
	for ; o.idx < len(o.src); o.idx++ {
		c := o.src[o.idx]
		if c == '"' && !isLiteral {
			isLiteral = true
		} else if c == '"' && isLiteral {
			isLiteral = false
		} else if isLiteral {
			continue
		}
		if c == '%' {
			hasRelFunc = true
		} else if hasRelFunc && c == '(' {
			hasRelFunc = false
			break
		} else if hasRelFunc {
			continue
		}
		if c == ' ' || c == '\t' || c == ',' || c == '(' || c == ')' {
			//if isDelim(o.src[o.idx]) {
			break
		}
	}
	// この関数に入った場合かならずtokenがある
	val := string(o.src[start:o.idx])
	typ := analyzeOperandType(val)
	return val, typ
}

func isRegister(val string) bool {
	return RegisterSet[val]
}

func isImmediate(value string) bool {
	// 正の整数リテラル（例: 42）
	integerPattern := `^\d+$`
	// 16進数リテラル（例: 0x2A）
	hexPattern := `^0x[0-9A-Fa-f]+$`

	// 数値の形式に合致するかを確認
	if matched, _ := regexp.MatchString(integerPattern, value); matched {
		// 数値として変換できるか確認
		_, err := strconv.Atoi(value)
		return err == nil
	}

	// 16進数リテラルとしての形式に合致するかを確認
	if matched, _ := regexp.MatchString(hexPattern, value); matched {
		// 16進数として変換できるか確認
		_, err := strconv.ParseInt(value[2:], 16, 64) // "0x" を除去して変換
		return err == nil
	}

	return false
}

func analyzeOperandType(operand string) OperandType {
	if isRegister(operand) {
		return REG
	} else if isImmediate(operand) {
		return IMM
	} else {
		return LAB
	}
}

// この関数に来る時点でラベルをオペランドにとることは確定している
func isValidRelFunc(opecode, relFunc string) bool {
	typ := OpecodeMap[opecode].opcTyp

	switch typ {
	case RType, IType, SType:
		if relFunc == "%lo" || relFunc == "%pcrel_lo" {
			return true
		}
		break

	case UType:
		if relFunc == "%hi" || relFunc == "%pcrel_hi" {
			return true
		}
		break

	default:
		break
	}
	return false
}

/*
命令形式ごとにオペランドが正しいか見る
*/
func (o *Operation) handleByOpType() error {
	oprTypIdx := 0
	o.skipUntilNextOperand()

	for !o.isEOF() && oprTypIdx < len(o.info.oprTyps) {
		val, typ := o.nextOperand()
		if o.info.oprTyps[oprTypIdx]&typ == 0 {
			return errors.New("illegal operand.")
		}
		// リロケーションファンクションの場合
		if typ == LAB && val[0] == '%' {
			if isValidRelFunc(o.Opecode(), val) {
				o.relFunc = val
				o.skipUntilNextOperand()
				continue
			}
			return errors.New("illegal operand.")
		}
		o.operands = append(o.operands, val)
		o.skipUntilNextOperand()
		oprTypIdx++
	}

	if oprTypIdx != len(o.info.oprTyps) {
		return errors.New("illegal operand.")
	} else if !o.isEOF() {
		return errors.New(fmt.Sprintf(ErrMsg, o.src[o.idx]))
	}
	return nil
}

func (s *Stmt) parseOperation(val string) error {
	op := Operation{
		opcode: val,
		info:   OpecodeMap[val],
		src:    s.src[s.idx:],
		idx:    0,
	}

	err := op.handleByOpType()
	if err != nil {
		return err
	}

	s.op = &op
	return nil
}
