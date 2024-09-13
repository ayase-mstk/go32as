package elf32

import "fmt"

const (
  // special section indexes
  SHN_UNDEF     = 0
  SHN_LORESERVE = 0xff00
  SHN_LOPROC    = 0xff00
  SHN_HIPROC    = 0xff1f
  SHN_LOOS      = 0xff20
  SHN_HIOS      = 0xff3f
  SHN_ABS       = 0xfff1
  SHN_COMMON    = 0xfff2
  SHN_XINDEX    = 0xffff
  SHN_HIRESERVE = 0xffff

  // symbol binding
  STB_LOCAL   = 0
  STB_GLOBAL  = 1
  STB_WEAK    = 2
  STB_LOOS    = 10
  STB_HIOS    = 12
  STB_LOPROC  = 13
  STB_HIPROC  = 15

  // symbol type
  STT_NOTYPE  = 0
  STT_OBJECT  = 1
  STT_FUNC    = 2
  STT_SECTION = 3
  STT_FILE    = 4
  STT_COMMON  = 5
  STT_TLS     = 6
  STT_LOOS    = 10
  STT_HIOS    = 12
  STT_LOPROC  = 13
  STT_HIPROC  = 15

  // symbol visibility
  STV_DEFAULT   = 0
  STV_INTERNAL  = 1
  STV_HIDDEN    = 2
  STV_PROTECTED = 3
)

type Elf32SymtblEntry struct {
  name  Elf32Word // symbol name index which stored in string table
  value Elf32Addr // symbol address value
  size  Elf32Word // symbol size
  info  uint8     // upper 4bit, lower 4bit symbol type - method or variable, local or global
  other uint8     // specifies a symbol's visibility.
  shndx Elf32Half // section header table index which symbol belongs to
}

type Symtbl struct {
  symtbls []Elf32SymtblEntry
  idx     map[string]int
}

var symbolInfoTypes = map[string]uint8{
    "@notype":     STT_NOTYPE,
    "@object":     STT_OBJECT,
    "@function":   STT_FUNC,
    "@section":    STT_SECTION,
    "@file":       STT_FILE,
    "@common":     STT_COMMON,
    "@tls_object": STT_TLS,
}

func (e *Elf32) initSymbolTables() {
  nilSymbol := Elf32SymtblEntry{
      name:    e.strtbl.resolveIndex(""),
      value:   0,
      size:    0,
      info:    0,
      other:   0,
      shndx:   SHN_UNDEF,
  }
  e.symtbl.addSymbol(nilSymbol, "")
}

func (s *Symtbl) addSymbol(sym Elf32SymtblEntry, name string) {
    // Initialize the idx map if it's not already initialized
    if s.idx == nil {
        s.idx = make(map[string]int)
    }

    // Add the new symbol entry to the symtbls slice
    s.symtbls = append(s.symtbls, sym)

    // Add the symbol name to idx with the index of the new entry
    s.idx[name] = len(s.symtbls) - 1
}

func createSymInfo(binding, typ byte) byte {
    return (binding << 4) | (typ & 0x0F)
}

func newSymbol(key string, name Elf32Word, value Elf32Addr, size Elf32Word, info byte, shndx Elf32Half) Elf32SymtblEntry {
  newEntry := Elf32SymtblEntry{
    name:   name,
    value:  value,
    size:   size,
    info:   info,
    shndx:  shndx,
  }

  return newEntry
}

func (s *Symtbl) setInfo(name string, info uint8) {
  id := s.idx[name]
  s.symtbls[id].info = info
}

func (s *Symtbl) setValue(name string, value Elf32Addr) {
  id := s.idx[name]
  s.symtbls[id].value = value
}

func (s *Symtbl) SymbolExists(name string) bool {
    _, exists := s.idx[name]
    return exists
}

func (s *Symtbl) printSymbolTable(strtbl Elf32Strtbl) {
  for _, sym := range s.symtbls {
    end := int(sym.name)
    for end < len(strtbl.data) && strtbl.data[end] != 0 {
      end++
    }
    fmt.Printf("name=%q\n", strtbl.data[sym.name:end])
    fmt.Println("value=", sym.value)
    fmt.Println("size=", sym.size)
    fmt.Println("info.binding=", sym.info >> 4)
    fmt.Println("info.type=", sym.info & 0x0F)
    fmt.Println("other=", sym.other)
    fmt.Println("shndx=", sym.shndx)
    fmt.Println("")
  }
}
