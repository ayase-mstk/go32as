package elf32

type Elf32Strtbl struct {
	data  []byte         // 連結された文字列のバッファ
	index map[string]int // 文字列からインデックスへのマッピング
}

func (st *Elf32Strtbl) exist(name string) bool {
	_, exists := st.index[name]
	return exists
}

func (st *Elf32Strtbl) resolveIndex(name string) Elf32Word {
	if st.exist(name) {
		return Elf32Word(st.index[name])
	}
	if st.index == nil {
		st.index = make(map[string]int)
	}

	// 現在の文字列テーブルの末尾インデックスを取得
	index := len(st.data)
	st.index[name] = index
	// シンボル名をnull終端付きで追加
	st.data = append(st.data, []byte(name)...)
	st.data = append(st.data, 0) // null終端
	return Elf32Word(index)
}
