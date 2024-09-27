    .section .data
message:
    .asciz "Hello World"  # 出力する文字列

    .section .text
    .globl _start

_start:
    # write(int fd, const void *buf, size_t count)
    # fd = 1 (stdout)
    # buf = address of message
    # count = length of the string

    # x10 (a0) = 1 (stdout)
    addi x10, x0, 1          # ファイルディスクリプタ(stdout)の値をレジスタx10に設定

    # x11 (a1) = address of the message
    # メモリのアドレスをロード (laの代わり)
    lui x11, %hi(message)    # messageの上位20ビットをロード
    addi x11, x11, %lo(message) # messageの下位12ビットを加算

    # x12 (a2) = length of the message
    addi x12, x0, 13         # 文字列の長さをx12に設定 ("Hello, World!\n"は13文字)

    # x17 (a7) = system call number (64 is sys_write)
    addi x17, x0, 64         # writeシステムコールの番号をx17に設定

    # system call
    ecall                    # システムコールを呼び出す

    # exit system call (終了)
    # x10 (a0) = 0 (exit status)
    addi x10, x0, 0          # 終了コード 0 (成功) をx10に設定

    # x17 (a7) = system call number (93 is sys_exit)
    addi x17, x0, 93         # exitシステムコールの番号をx17に設定
    ecall                    # システムコールを呼び出す

