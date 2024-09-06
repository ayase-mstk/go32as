package parsetest

import (
	"fmt"
	"github.com/ayase-mstk/go32as/src/parse"
	"testing"
)

type parseDirectiveTestStruct struct {
	expectedVal   string
	expectedLabel string
}

func expectSameDirective(t *testing.T, got parse.Stmt, want []parseDirectiveTestStruct) {
	expectSameSize(t, len(got.Dir().Args())+1, len(want))

	for i, tt := range want {
		actual := func() string {
			if i == 0 {
				return got.Dir().Name()
			}
			return got.Dir().Args()[i-1]
		}()

		if actual != tt.expectedVal {
			t.Fatalf("test[%d] - section val wrong. got=%q, expected=%q",
				i, actual, tt.expectedVal)
		}

		if got.LSymbol() != tt.expectedLabel {
			t.Fatalf("test[%d] - label wrong. got=%q, expected=%q",
				i, got.LSymbol(), tt.expectedLabel)
		}
	}
}

func TestParseDirectiveAlign(t *testing.T) {
	input := []rune(" .align 4")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".align",
			expectedLabel: "",
		},
		{
			expectedVal:   "4",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveFile(t *testing.T) {
	input := []rune(" .file \"main.s\"")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".file",
			expectedLabel: "",
		},
		{
			expectedVal:   "\"main.s\"",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveGlobl(t *testing.T) {
	input := []rune("  .globl _start")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".globl",
			expectedLabel: "",
		},
		{
			expectedVal:   "_start",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveLocal(t *testing.T) {
	input := []rune("  .local myVar")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".local",
			expectedLabel: "",
		},
		{
			expectedVal:   "myVar",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveComm(t *testing.T) {
	input := []rune("  .comm myArray, 128, 4")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".comm",
			expectedLabel: "",
		},
		{
			expectedVal:   "myArray",
			expectedLabel: "",
		},
		{
			expectedVal:   "128",
			expectedLabel: "",
		},
		{
			expectedVal:   "4",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveIdent(t *testing.T) {
	input := []rune("  .ident \"GCC: (GNU) 10.2.0\"")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".ident",
			expectedLabel: "",
		},
		{
			expectedVal:   "\"GCC: (GNU) 10.2.0\"",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveSectionText(t *testing.T) {
	input := []rune("  .section .text")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".section",
			expectedLabel: "",
		},
		{
			expectedVal:   ".text",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveSectionData(t *testing.T) {
	input := []rune("  .section .data")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".section",
			expectedLabel: "",
		},
		{
			expectedVal:   ".data",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveSectionRodata(t *testing.T) {
	input := []rune("  .section .rodata")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".section",
			expectedLabel: "",
		},
		{
			expectedVal:   ".rodata",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveSectionBss(t *testing.T) {
	input := []rune("  .section .bss")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".section",
			expectedLabel: "",
		},
		{
			expectedVal:   ".bss",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveSize(t *testing.T) {
	input := []rune("  .size myFunc, 32")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".size",
			expectedLabel: "",
		},
		{
			expectedVal:   "myFunc",
			expectedLabel: "",
		},
		{
			expectedVal:   "32",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveString(t *testing.T) {
	input := []rune("  .string \"Hello, Null-Terminated!\"")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".string",
			expectedLabel: "",
		},
		{
			expectedVal:   "\"Hello, Null-Terminated!\"",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveAsciz(t *testing.T) {
	input := []rune("  .asciz \"Hello, Nice to meet you!\"")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".asciz",
			expectedLabel: "",
		},
		{
			expectedVal:   "\"Hello, Nice to meet you!\"",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveEqu(t *testing.T) {
	input := []rune("  .equ myConst, 42")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".equ",
			expectedLabel: "",
		},
		{
			expectedVal:   "myConst",
			expectedLabel: "",
		},
		{
			expectedVal:   "42",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveMacro(t *testing.T) {
	input := []rune("  .macro myMacro")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".macro",
			expectedLabel: "",
		},
		{
			expectedVal:   "myMacro",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveEndm(t *testing.T) {
	input := []rune("  .endm")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".endm",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveType(t *testing.T) {
	input := []rune("  .type myFunc, @function")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".type",
			expectedLabel: "",
		},
		{
			expectedVal:   "myFunc",
			expectedLabel: "",
		},
		{
			expectedVal:   "@function",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveByte(t *testing.T) {
	input := []rune("  .byte 0x12")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".byte",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x12",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirective2Byte(t *testing.T) {
	input := []rune("  .2byte 0x1234")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".2byte",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x1234",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveHalf(t *testing.T) {
	input := []rune("  .half 0x1234")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".half",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x1234",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirective4Byte(t *testing.T) {
	input := []rune("  .4byte 0x12345678")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".4byte",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x12345678",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveWord(t *testing.T) {
	input := []rune("  .word 0x12345678")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".word",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x12345678",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveLong(t *testing.T) {
	input := []rune("  .long 0x12345678")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".long",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x12345678",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveZero(t *testing.T) {
	input := []rune("  .zero 4")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".zero",
			expectedLabel: "",
		},
		{
			expectedVal:   "4",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

func TestParseDirectiveAttribute(t *testing.T) {
	input := []rune("  .attribute foo, 0x10")
	stmt, err := parse.ParseLine(input, 1)
	if err != nil {
		t.Fatalf("test - parse failed:\n%q", err.Error())
	}

	tests := []parseDirectiveTestStruct{
		{
			expectedVal:   ".attribute",
			expectedLabel: "",
		},
		{
			expectedVal:   "foo",
			expectedLabel: "",
		},
		{
			expectedVal:   "0x10",
			expectedLabel: "",
		},
	}

	expectSameDirective(t, stmt, tests)
}

/*
======================================
=========== Error Test ===============
======================================
*/

// Error Message
const (
	BadExpression   = "bad or irreducible absolute expression."
	MissingArgument = "missing argument."
)

func TestParseDirectiveErrorAlign(t *testing.T) {
	input := []rune("  .align myFunc")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'm'))
}

func TestParseDirectiveErrorFile(t *testing.T) {
	input := []rune("  .file")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorGlobl(t *testing.T) {
	input := []rune("  .globl")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorLocal(t *testing.T) {
	input := []rune("  .local myFunc, 1")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '1'))
}

func TestParseDirectiveErrorComm(t *testing.T) {
	input := []rune("  .comm")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorCommon(t *testing.T) {
	input := []rune("  .common 42, 1024, 4")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '4'))
}

func TestParseDirectiveErrorIdent(t *testing.T) {
	input := []rune("  .ident")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorSection(t *testing.T) {
	input := []rune("  .section text, 'ax'")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '\''))
}

func TestParseDirectiveErrorSize(t *testing.T) {
	input := []rune("  .size 0x1000")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '0'))
}

func TestParseDirectiveErrorText(t *testing.T) {
	input := []rune("  .text 1")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}
	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '1'))
}

func TestParseDirectiveErrorData(t *testing.T) {
	input := []rune("  .data myFunc")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}
	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'm'))
}

func TestParseDirectiveErrorRoData(t *testing.T) {
	input := []rune("  .rodata abc")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}
	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'a'))
}

func TestParseDirectiveErrorBss(t *testing.T) {
	input := []rune("  .bss .section")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}
	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '.'))
}

func TestParseDirectiveErrorString(t *testing.T) {
	input := []rune("  .string")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorAsciz(t *testing.T) {
	input := []rune("  .asciz")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorEqu(t *testing.T) {
	input := []rune("  .equ")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorMacro(t *testing.T) {
	input := []rune("  .macro")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}

func TestParseDirectiveErrorEndm(t *testing.T) {
	input := []rune("  .endm func")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'f'))
}

func TestParseDirectiveErrorType(t *testing.T) {
	input := []rune("  .type myFunc, 12")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '1'))
}

func TestParseDirectiveErrorByte(t *testing.T) {
	input := []rune("  .byte extra")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'e'))
}

func TestParseDirectiveErrorByte2(t *testing.T) {
	input := []rune("  .2byte abc")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'a'))
}

func TestParseDirectiveErrorHalf(t *testing.T) {
	input := []rune("  .half extra")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'e'))
}

func TestParseDirectiveErrorShort(t *testing.T) {
	input := []rune("  .short string")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 's'))
}

func TestParseDirectiveErrorByte4(t *testing.T) {
	input := []rune("  .4byte extra")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'e'))
}

func TestParseDirectiveErrorWord(t *testing.T) {
	input := []rune("  .word 0x12345678@")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, '0'))
}

func TestParseDirectiveErrorLong(t *testing.T) {
	input := []rune("  .long abc")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'a'))
}

func TestParseDirectiveErrorZero(t *testing.T) {
	input := []rune("  .zero extra")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), fmt.Sprintf(UnrecognizedError, 'e'))
}

func TestParseDirectiveErrorAttribute(t *testing.T) {
	input := []rune("  .attribute")
	_, err := parse.ParseLine(input, 1)
	if err == nil {
		t.Fatalf("test - parse have to be fail.")
	}

	expectErrorMessage(t, err.Error(), MissingArgument)
}
