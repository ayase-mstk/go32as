package main

import (
	"fmt"
	"os"

	"github.com/ayase-mstk/go32as/src/elf32"
	"github.com/ayase-mstk/go32as/src/parse"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "invalid num of arguments.")
		os.Exit(0)
	}

	stmts, err := parse.ParseFile(os.Args[1])
	if err != nil {
		fmt.Printf("%s: Assembler messages:\n", os.Args[1])
		fmt.Println(err.Error())
		os.Exit(0)
	}

	e, err := elf32.PrepareElf32Tables(stmts)
	if err != nil {
		fmt.Printf("%s: Assembler messages:\n", os.Args[1])
		fmt.Println(os.Args[1], ":", err.Error())
		os.Exit(0)
	}
	//e.PrintAll()
	e.WriteToFile()
}
