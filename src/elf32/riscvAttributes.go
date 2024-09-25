package elf32

import (
	"unsafe"
)

// uleb128 (Unsigned LEB128) は可変長エンコーディングをサポートするため、カスタムの型や関数で表現します。
type ULEB128 uint64

// ULEB128 のサイズを計算するためのヘルパー関数
func uleb128Size(value ULEB128) Elf32Word {
	size := 0
	for value >= 0x80 {
		size++
		value >>= 7
	}
	return Elf32Word(size + 1) // 最後の1バイトを加算
}

// .riscv.attributes section.
type Elf32Attributes struct {
	FormatVersion  byte            // The format version (e.g., 'A')
	VendorSections []VendorSection // Vendor-specific subsections
}

// VendorSection represents a vendor-specific subsection.
type VendorSection struct {
	Length         Elf32Word       // Length of the subsection
	VendorName     string          // Vendor name (NTBS, null-terminated byte string)
	SubSubSections []SubSubSection // One or more sub-sub-sections
}

// SubSubSection represents a sub-sub-section with attributes.
type SubSubSection struct {
	Tag        ULEB128     // The tag for the sub-sub-section
	Length     Elf32Word   // Length of the sub-sub-section
	Attributes []Attribute // Tag-value pairs for this section
}

// Attribute represents a single attribute with a tag-value pair.
type Attribute struct {
	Tag   ULEB128     // The attribute's tag
	Value interface{} // The value can be either ULEB128 or a string (NTBS)
}

// Predefined attribute types (as examples).
type RISCVAttributes struct {
	StackAlign      ULEB128 // Tag 4, uleb128, stack alignment (e.g., 16 bytes for RV32I)
	Arch            string  // Tag 5, NTBS, architecture string
	UnalignedAccess ULEB128 // Tag 6, ulem2p0b128, unaligned access (e.g., 0 or 1)
	AtomicABI       ULEB128 // Tag 14, uleb128, atomic ABI version
	X3RegUsage      ULEB128 // Tag 16, uleb128, usage of x3/gp register
}

// Helper function to create a new attribute (for example, for stack alignment).
func NewAttribute(tag ULEB128, value interface{}) Attribute {
	return Attribute{Tag: tag, Value: value}
}

// CalculateLength はサブサブセクションの正確な Length を計算します
func (s *SubSubSection) CalculateLength() Elf32Word {
	length := Elf32Word(0)

	// タグのサイズ (ULEB128) を計算
	length += uleb128Size(s.Tag)

	// 各属性のサイズを計算
	for _, attr := range s.Attributes {
		// タグのサイズ (ULEB128)
		length += uleb128Size(attr.Tag)

		// 値のサイズ (ULEB128 または NTBS)
		switch v := attr.Value.(type) {
		case ULEB128:
			length += uleb128Size(v)
		case string:
			length += Elf32Word(len(v)) + 1 // NTBS の長さ + 終端の null バイト
		}
	}

	return length
}

// Helper function to create a new vendor section.
func NewVendorSection(name string, attributes []Attribute) VendorSection {
	subSection := SubSubSection{
		Tag: 1, // Tag_file, relating to the whole file
		//Length:     uint32(len(attributes) * 8), // Placeholder for sub-sub-section length
		Attributes: attributes,
	}
	subSection.Length = subSection.CalculateLength()
	return VendorSection{
		Length:         Elf32Word(len(attributes)*8 + len(name) + 1),
		VendorName:     name,
		SubSubSections: []SubSubSection{subSection},
	}
}

// Example usage
func (e *Elf32) initAttributes() {
	// Define some example attributes
	attrs := []Attribute{
		NewAttribute(4, ULEB128(16)), // Stack alignment: 16 bytes
		NewAttribute(5, "rv32i2p1"),  // Architecture: RV32I
		NewAttribute(6, ULEB128(0)),  // Unaligned access: not allowed
		NewAttribute(14, ULEB128(0)), // Atomic ABI: no
		NewAttribute(16, ULEB128(0)), // x3 register usage: default usage
	}

	// Create a vendor section for "riscv"
	riscvVendor := NewVendorSection("riscv", attrs)

	// Create and return the full attributes section
	e.attr = Elf32Attributes{
		FormatVersion:  'A', // Format version 'A'
		VendorSections: []VendorSection{riscvVendor},
	}
}

// CalculateSize calculates the size of an Attribute in bytes.
func (attr *Attribute) CalculateSize() Elf32Word {
	size := Elf32Word(unsafe.Sizeof(attr.Tag)) // タグのサイズ
	switch v := attr.Value.(type) {
	case ULEB128:
		size += uleb128Size(v) // ULEB128 の値サイズ
	case string:
		size += Elf32Word(len(v)) + 1 // 文字列の長さ + 終端の null バイト
	}
	return size
}

// CalculateSize calculates the size of a SubSubSection in bytes.
func (sss *SubSubSection) CalculateSize() Elf32Word {
	size := Elf32Word(unsafe.Sizeof(sss.Tag))    // サブサブセクションタグのサイズ
	size += Elf32Word(unsafe.Sizeof(sss.Length)) // Length フィールドのサイズ (4 bytes)
	for _, attr := range sss.Attributes {
		size += attr.CalculateSize() // 各属性のサイズ
	}
	return size
}

// CalculateSize calculates the size of a VendorSection in bytes.
func (vs *VendorSection) CalculateSize() Elf32Word {
	size := Elf32Word(unsafe.Sizeof(vs.Length)) // Length フィールドのサイズ (4 bytes)
	size += Elf32Word(len(vs.VendorName)) + 1   // Vendor name のサイズ (NTBS)

	for _, sub := range vs.SubSubSections {
		size += sub.CalculateSize() // 各サブサブセクションのサイズ
	}
	return size
}

// CalculateSize calculates the total size of the Elf32Attributes section in bytes.
func (as *Elf32Attributes) CalculateSize() Elf32Word {
	size := Elf32Word(unsafe.Sizeof(as.FormatVersion)) // フォーマットバージョン (1 byte)
	for _, vendor := range as.VendorSections {
		size += vendor.CalculateSize() // 各ベンダーセクションのサイズ
	}
	return size
}
