# RV32I Assembler

### 概要
このプロジェクトは、RISC-V（RV32I）命令セットアーキテクチャ用の自作アセンブラです。<br>
RV32I基本命令セットをサポートし、アセンブリ言語で書かれたコード(*.s)をELFフォーマットに従って機械語に変換し、オブジェクトファイル(*.o)を作成します。

### サポート事項
RV32I基本命令セットの完全サポート<br>
多くのgABIディレクティブのサポート<br>
シンボルテーブルとラベル解決機能<br>
エラー検出とエラーメッセージの出力<br>

### 非サポート事項
疑似命令のサポート<br>
一部重要でないディレクティブ<br>
詳細なエラー検出<br>

### 事前準備
[riscv-gnu-toolchain](https://github.com/riscv-collab/riscv-gnu-toolchain), [spike](https://github.com/riscv-software-src/riscv-isa-sim), [pk](https://github.com/riscv-software-src/riscv-pk)のinstall<br>
こちらの記事を参考に環境構築しました。（https://zenn.dev/ohno418/articles/5f6d5e01dc4981）
```
// path
$ export RISCV=/path/to/riscv/tools
$ export PATH=$RISCV/bin:$PATH

// riscv-gnu-toolchain for rv32i
$ ./configure --prefix=$RISCV --with-arch=rv32i --with-abi=ilp32
$ make
$ make install

// spike
$ mkdir build && cd build
$ sudo ../configure --prefix=$RISCV --with-isa=rv32i
$ sudo make
$ sudo make install

// pk
$ mkdir build && cd build
$ ../configure --prefix=$RISCV --host=riscv32-unknown-linux-gnu --with-arch=rv32i --with-abi=ilp32
$ sudo make
$ sudo make install
```

### 使い方
```
make
./rv32i-as sample/helloworld.s
path/to/riscv32-unknown-linux-gnu-gcc -static -nostartfiles output.o -o a.out
path/to/spike path/to/pk a.out
```

### 命令
RV32I命令セットをすべてサポートしています。<br>
主な命令には以下が含まれます：
```
算術演算命令: ADD, SUB, ADDI など
論理演算命令: AND, OR, XOR など
分岐命令: BEQ, BNE, BLT など
ロード/ストア命令: LW, SW など
ジャンプ命令: JAL, JALR
```

完全な命令リストについては、[RISCV REFERENCE CARD](https://www.cs.sfu.ca/~ashriram/Courses/CS295/assets/notebooks/RISCV/RISCV_CARD.pdf)を参照してください。


### ディレクティブ
[riscv-asm-manual](https://github.com/riscv-non-isa/riscv-asm-manual/blob/main/src/asm-manual.adoc#pseudo-ops)に記載されているほとんどの32bit向けディレクティブをサポートしています。<br>
主なディレクティブには以下が含まれます：
```
シンボル関連：　.local, .globl, size, .type
セクション関連： .section, .text, .data, .bss, .rodata
データ関連：　.string, .word
```

### その他参考文献
[gABI ELF format 仕様書](https://www.sco.com/developers/gabi/latest/contents.html)<br>
[riscv elf format](https://github.com/riscv-non-isa/riscv-elf-psabi-doc/blob/master/riscv-elf.adoc#elf-object-files)<br>
