package elf32

type Elf32Strtbl struct {
    data   []byte            // 連結された文字列のバッファ
    index  map[string]int    // 文字列からインデックスへのマッピング
}

func (st *Elf32Strtbl) resolveIndex(name string) Elf32Word {
    // 既に存在していれば追加する
    if _, exist := st.index[name]; exist {
      return Elf32Word(st.index[name])
    }

    // 現在の文字列テーブルの末尾インデックスを取得
    index := len(st.data)
    // シンボル名をnull終端付きで追加
    st.data = append(st.data, []byte(name)...)
    st.data = append(st.data, 0) // null終端
    return Elf32Word(index)
}
