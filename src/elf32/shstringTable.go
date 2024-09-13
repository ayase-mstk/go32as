package elf32

type Elf32Shstrtbl struct {
    data   []byte            // 連結された文字列のバッファ
    idx    map[string]int    // 文字列からインデックスへのマッピング
}

func (s *Elf32Shstrtbl) resolveIndex(name string) Elf32Word {
    // 既に存在していれば追加する
    if _, exist := s.idx[name]; exist {
      return Elf32Word(s.idx[name])
    }

    // 現在の文字列テーブルの末尾インデックスを取得
    idx := len(s.data)
    // シンボル名をnull終端付きで追加
    s.data = append(s.data, []byte(name)...)
    s.data = append(s.data, 0) // null終端
    return Elf32Word(idx)
}
