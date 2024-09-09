package elf32

import (
  "unsafe"
)

const (
	// ELF識別子のサイズ
	EiNident = 16

	// ELFファイルタイプ
	ETNone = 0 // 未定義
	ETRel  = 1 // 再配置可能ファイル
	ETExec = 2 // 実行可能ファイル
	ETDyn  = 3 // 共有オブジェクトファイル
	ETCore = 4 // コアファイル

	// マシンアーキテクチャ
	EMNone  = 0   // 未定義
	EMRiscv = 243 // RISC-V

	// ELFバージョン
	EVNone    = 0 // 無効
	EVCurrent = 1 // 現行バージョン
)

// ELF32ヘッダー構造体
type Elf32Ehdr struct {
	EIdent     [EiNident]byte // ELF識別子
	EType      Elf32Half      // ELFファイルのタイプ
	EMachine   Elf32Half      // マシンアーキテクチャ
	EVersion   Elf32Word      // ELFバージョン
	EEntry     Elf32Addr      // エントリポイントアドレス
	EPhoff     Elf32Off       // プログラムヘッダーのオフセット
	EShoff     Elf32Off       // セクションヘッダーのオフセット
	EFlags     Elf32Word      // プロセッサ特有のフラグ
	EEhsize    Elf32Half      // ELFヘッダーのサイズ
	EPhentsize Elf32Half      // プログラムヘッダーエントリのサイズ
	EPhnum     Elf32Half      // プログラムヘッダーのエントリ数
	EShentsize Elf32Half      // セクションヘッダーエントリのサイズ
	EShnum     Elf32Half      // セクションヘッダーのエントリ数
	EShstrndx  Elf32Half      // セクション名文字列テーブルのインデックス
}

func (e *Elf32) initHeader() {
	// initialize
	e.ehdr = Elf32Ehdr{
		EIdent: [EiNident]byte{
			0x7f, 0x45, 0x4c, 0x46, // 0x7F followed by "ELF"(45 4c 46) in ASCII;
			0x01,                                     // EI_CLASS:1=32-bit
			0x01,                                     // EI_DATA:1=little endian
			0x01,                                     // EI_VERSION:1=the original and current version of ELF.
			0x00,                                     // EI_OSABI: 0=System V
			0x00,                                     // EI_ABIVERSION:
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // EI_PAD: always zero.
		},
		EType:      ETRel,
		EMachine:   EMRiscv,
		EVersion:   EVCurrent,
		EEntry:     0, // always zero in EVRel
		EPhoff:     0, // always zero in EVRel
		EShoff:     0, // entrypoint of section header
		EFlags:     0, // always zero in RV32I
		EEhsize:    Elf32Half(unsafe.Sizeof(Elf32Ehdr{})),
		EPhentsize: 0,
		EPhnum:     0,
		EShentsize: Elf32Half(unsafe.Sizeof(Elf32Shdr{})),
		EShnum:     0,
		EShstrndx:  0,
	}
}
