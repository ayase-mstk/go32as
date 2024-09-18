package elf32

func calcOffset(off Elf32Off, size Elf32Word) Elf32Off {
	return Elf32Off(int(off) + int(size))
}

func (e *Elf32) ResolveSectionRayout() {
	var lastOffset Elf32Off = 0x34
	//for _, shdr := range e.shdr.shdrs {
	//  // size
	//  name := e.shstrtbl.data[shdr.ShName]
	//  // offset
	//  lastOffset += size
	//}

	// .text section
	e.shdr.setSize(".text", Elf32Word(e.sections.entry[".text"].off))
	lastOffset += Elf32Off(e.shdr.getSize(".text"))

	// .data section
	e.shdr.setOffset(".data", lastOffset)
	e.shdr.setSize(".data", Elf32Word(e.sections.entry[".data"].off))
	lastOffset += Elf32Off(e.shdr.getSize(".data"))

	// .bss section
	e.shdr.setOffset(".bss", lastOffset)
	e.shdr.setSize(".bss", Elf32Word(e.sections.entry[".bss"].off))
	lastOffset += Elf32Off(e.shdr.getSize(".bss"))

	// .rodata section
	if _, exist := e.shdr.shndx[".rodata"]; exist {
		e.shdr.setOffset(".rodata", lastOffset)
		e.shdr.setSize(".rodata", Elf32Word(e.sections.entry[".rodata"].off))
		lastOffset += Elf32Off(e.shdr.getSize(".rodata"))
	}

	// .riscv.attributes section
	e.shdr.setOffset(".riscv.attributes", lastOffset)
	e.shdr.setSize(".riscv.attributes", 0x5f)
	lastOffset += Elf32Off(e.shdr.getSize(".riscv.attributes"))

	// .symtab section
	e.shdr.setOffset(".symtab", lastOffset)
	e.shdr.setSize(".symtab", Elf32Word(len(e.symtbl.symtbls))*e.shdr.getEntsize(".symtab"))
	e.shdr.setLink(".symtab", Elf32Word(e.shdr.shndx[".strtab"]))
	e.shdr.setInfo(".symtab", Elf32Word(e.symtbl.calcLastLocalSymIdx()+1))
	e.shdr.setEntsize(".symtab", 0x10)
	lastOffset += Elf32Off(e.shdr.getSize(".symtab"))

	// .strtab section
	e.shdr.setOffset(".strtab", lastOffset)
	e.shdr.setOffset(".strtab", calcOffset(e.shdr.getOffset(".symtab"), e.shdr.getSize(".symtab")))
	e.shdr.setSize(".strtab", Elf32Word(len(e.strtbl.data)))
	lastOffset += Elf32Off(e.shdr.getSize(".strtab"))

	// .strtab section
	e.shdr.setOffset(".shstrtab", lastOffset)
	e.shdr.setSize(".shstrtab", Elf32Word(len(e.shstrtbl.data)))
	lastOffset += Elf32Off(e.shdr.getSize(".shstrtab"))

	// .rela.text section
	if _, exist := e.shdr.shndx[".rela.text"]; exist {
		e.shdr.setOffset(".rela.text", lastOffset)
		e.shdr.setSize(".rela.text", Elf32Word(len(e.rela.entry))*e.shdr.getEntsize(".rela.text"))
		e.shdr.setLink(".rela.text", Elf32Word(e.shdr.shndx[".symtab"]))
		e.shdr.setInfo(".symtab", Elf32Word(e.shdr.shndx[".text"]))
		lastOffset += Elf32Off(e.shdr.getSize(".rela.text"))
	}
}
