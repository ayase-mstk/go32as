package elf32

type Elf32Shstrtbl struct {
	data []byte         // 連結された文字列のバッファ
	idx  map[string]int // 文字列からインデックスへのマッピング
}

func (s *Elf32Shstrtbl) exist(name string) bool {
	_, exists := s.idx[name]
	return exists
}

func (s *Elf32Shstrtbl) resolveIndex(name string) Elf32Word {
	// 既に存在していれば追加する
	if s.exist(name) {
		return Elf32Word(s.idx[name])
	}
	if s.idx == nil {
		s.idx = make(map[string]int)
	}

	// 現在の文字列テーブルの末尾インデックスを取得
	idx := len(s.data)
	s.idx[name] = idx
	// シンボル名をnull終端付きで追加
	s.data = append(s.data, []byte(name)...)
	s.data = append(s.data, 0) // null終端
	return Elf32Word(idx)
}
