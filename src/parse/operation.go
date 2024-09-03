package parse

// import "strings"

const (
  // R format
  ADD    = "add"
  SUB    = "sub"
  XOR    = "xor"
  OR     = "or"
  AND    = "and"
  SLL    = "sll"
  SRL    = "srl"
  SRA    = "sra"
  SLT    = "slt"
  SLTU   = "sltu"

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
  ECALL  = "ecall"
  EBREAK = "ebreak"

  // S format
  SB     = "sb"
  SH     = "sh"
  SW     = "sw"

  // B format
  BEQ    = "beq"
  BNE    = "bne"
  BLT    = "bge"
  BLTU   = "bltu"
  BGEU   = "bgeu"

  // J format
  JAL    = "jal"

  // U format
  LUI    = "lui"
  AUIPC  = "auipc"
)

var OpecodeSet = map[string]struct{}{
  ADD:    {},
  SUB:    {},
  XOR:    {},
  OR:     {},
  AND:    {},
  SLL:    {},
  SRL:    {},
  SRA:    {},
  SLT:    {},
  SLTU:   {},
  ADDI:   {},
  WORI:   {},
  ORI:    {},
  ANDI:   {},
  SLLI:   {},
  SRLI:   {},
  SRAI:   {},
  SLTI:   {},
  SLTIU:  {},
  LB:     {},
  LH:     {},
  LW:     {},
  LBU:    {},
  LHU:    {},
  ECALL:  {},
  EBREAK: {},
  SB:     {},
  SH:     {},
  SW:     {},
  BEQ:    {},
  BNE:    {},
  BLT:    {},
  BLTU:   {},
  BGEU:   {},
  JAL:    {},
  LUI:    {},
  AUIPC:  {},
}

type OpecodeType int

const (
  RType OpecodeType = iota
  IType
  SType
  BTYpe
  JType
  UType
)

type Opecode struct {
  Args  ArgFlag
  Type  OpecodeType
}


