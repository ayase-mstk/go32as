// parse/operation_parser_test.go

package parsetest

import (
	"github.com/ayase-mstk/go32as/src/parse"
	"testing"
)

type parseOperationTestStruct struct {
	expectedOpcType parse.OpecodeType
	expectedOprType parse.OperandType
	expectedVal     string
	expectedRow     int
}

func expectSameSize(t *testing.T, got, want int) {
	if got != want {
		t.Fatalf("token num wrong. actual = %d, expected = %d",
			got, want)
	}
}

func expectSameOperation(t *testing.T, got parse.Stmt, want []parseOperationTestStruct) {
	expectSameSize(t, len(got.Op().Operands())+1, len(want))

	for i, tt := range want {

		if i == 0 {
			if got.Op().OpcType() != tt.expectedOpcType {
				t.Fatalf("test[%d] - OpecodeType wrong. got=%q, expected=%q",
					i, got.Op().OpcType(), tt.expectedOpcType)
			}
			if got.Op().Opecode() != tt.expectedVal {
				t.Fatalf("test[%d] - opecode wrong. got=%q, expected=%q",
					i, got.Op().Opecode(), tt.expectedVal)
			}
		} else {
			if got.Op().OprType()[i-1]&tt.expectedOprType == 0 {
				t.Fatalf("test[%d] - OperandType wrong. got=%q, expected=%q",
					i, got.Op().OprType()[i-1], tt.expectedOprType)
			}
			if got.Op().Operands()[i-1] != tt.expectedVal {
				t.Fatalf("test[%d] - operand wrong. got=%q, expected=%q",
					i, got.Op().Operands()[i-1], tt.expectedVal)
			}
		}

		if got.Row() != tt.expectedRow {
			t.Fatalf("test[%d] - row wrong. got=%d, expected=%d",
				i, got.Row(), tt.expectedRow)
		}
	}
}

func TestParseOperation1(t *testing.T) {
	input := []rune("addi a0, a1, 42")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.IType,
			expectedVal:     "addi",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "a0",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "a1",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.IMM,
			expectedVal:     "42",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

func TestParseOperation2(t *testing.T) {
	input := []rune("    sw ra, 32(sp)")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.SType,
			expectedVal:     "sw",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "ra",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.IMM,
			expectedVal:     "32",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "sp",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

func TestParseOperation3(t *testing.T) {
	input := []rune("    jalr x0,x1,0")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.IType,
			expectedVal:     "jalr",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "x0",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "x1",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.IMM,
			expectedVal:     "0",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

func TestParseOperation4(t *testing.T) {
	input := []rune("    ecall")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.IType,
			expectedVal:     "ecall",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

func TestParseOperation5(t *testing.T) {
	input := []rune("   jal a0,func")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.JType,
			expectedVal:     "jal",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "a0",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.LAB,
			expectedVal:     "func",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

// まだ失敗する
func TestParseOperationRelocation(t *testing.T) {
	input := []rune("    lui a5,%hi(.LC0)")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseOperationTestStruct{
		{
			expectedOpcType: parse.JType,
			expectedVal:     "lui",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.REG,
			expectedVal:     "a5",
			expectedRow:     1,
		},
		{
			expectedOprType: parse.LAB,
			expectedVal:     "%h(.LC0)",
			expectedRow:     1,
		},
	}

	expectSameOperation(t, stmt, tests)
}

/*
=====================================
=========== Error Test ==============
=====================================
*/

const (
	OperandErr = "illegal operand."
)

func TestParseOperationError1(t *testing.T) {
	input := []rune("add a0, a1, 42")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}

func TestParseOperationError2(t *testing.T) {
	input := []rune("addi a0, a1, a2")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}

func TestParseOperationError3(t *testing.T) {
	input := []rune("sw ")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}

func TestParseOperationError4(t *testing.T) {
	input := []rune("lw s0 42(sp)")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}

func TestParseOperationError5(t *testing.T) {
	input := []rune("sw s0,42,sp")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}

func TestParseOperationError6(t *testing.T) {
	input := []rune("sb s0(42),sp")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), OperandErr)
}
