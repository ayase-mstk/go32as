package elf32

import "github.com/ayase-mstk/go32as/src/parse"

// 型にしたほうが関数呼び出ししやすい
type Elf32Sections struct {
	entry map[string]Section
}

type Section struct {
	stmts []parse.Stmt
	off   Elf32Addr // section 内のオフセット
}

func (s *Elf32Sections) resolveOffset(name string) Elf32Addr {
	return s.entry[name].off
}

func (s *Elf32Sections) advanceOffset(name string, off Elf32Addr) {
	if s.entry == nil {
		s.entry = make(map[string]Section)
	}
	section := s.entry[name]
	section.off += off
	s.entry[name] = section
}

func (s *Elf32Sections) exist(name string) bool {
	_, exists := s.entry[name]
	return exists
}
