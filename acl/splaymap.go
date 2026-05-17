package acl

import "cmp"

// Splaymapはキー順に並んだmap型のスプレー木
type Splaymap[K cmp.Ordered, V any] struct {
	root *splaynode[K, V]
}

// NewSplaymapは新しいスプレー木を作成
func NewSplaymap[K cmp.Ordered, V any]() *Splaymap[K, V] {
	return &Splaymap[K, V]{}
}

// Sizeは要素数を返す
func (st *Splaymap[K, V]) Size() int {
	if st.root == nil {
		return 0
	}
	return st.root.size
}

// IsEmptyは空かどうかを返す
func (st *Splaymap[K, V]) IsEmpty() bool {
	return st.Size() == 0
}

// Getはキーに対応する値と存在フラグを返す
func (st *Splaymap[K, V]) Get(key K) (value V, found bool) {
	if st.root == nil {
		return
	}

	st.root, found = st.root.has(key)
	if found {
		value = st.root.value
	}
	return
}

// Setはキーと値のペアを設定する。
// 既存キーなら値を更新、なければ新規挿入。
func (st *Splaymap[K, V]) Set(key K, value V) {
	L, R := split(st.root, key, SplitLE_GT)
	if L != nil {
		L = L.splayMax()
		if L.key == key {
			// Update
			L.value = value
			st.root = merge(L, R)
			return
		}
	}

	// 新ノード追加
	newNode := &splaynode[K, V]{key: key, value: value, size: 1}
	st.root = merge(merge(L, newNode), R)
}

// Deleteはキーで要素を削除し、削除できたかを返す
func (st *Splaymap[K, V]) Delete(key K) (deleted bool) {
	L, R := split(st.root, key, SplitLE_GT)
	if L != nil {
		L = L.splayMax()
	}

	if L == nil || L.key != key {
		// 対象なし。削除失敗
		st.root = merge(L, R)
		return
	}

	// 削除成功
	deleted = true
	left, _ := L.cutLeft()
	st.root = merge(left, R)
	return
}

// Atは昇順でk番目(0-indexed)の要素を返す。
// 範囲外ならok=false。
func (st *Splaymap[K, V]) At(k int) (key K, value V, ok bool) {
	if st.root == nil {
		return
	}

	if k < 0 || k >= st.Size() {
		return
	}

	st.root = st.root.kth(k)
	return st.root.key, st.root.value, true
}

// FirstGtはkeyより大きい最小要素のindexを返す。
// 存在しなければSize()を返す。要素はAt(FirstGt(key))で取得。
func (st *Splaymap[K, V]) FirstGt(key K) (idx int) {
	if st.root == nil {
		return
	}
	L, R := split(st.root, key, SplitLE_GT)
	if L != nil {
		idx = L.size
	}
	st.root = merge(L, R)
	return
}

// FirstGeはkey以上の最小要素のindexを返す。
// 存在しなければSize()を返す。要素はAt(FirstGe(key))で取得。
func (st *Splaymap[K, V]) FirstGe(key K) (idx int) {
	if st.root == nil {
		return
	}
	L, R := split(st.root, key, SplitLT_GE)
	if L != nil {
		idx = L.size
	}
	st.root = merge(L, R)
	return
}

// LastLtはkeyより小さい最大要素のindexを返す。
// 存在しなければ-1を返す。要素はAt(LastLt(key))で取得。
func (st *Splaymap[K, V]) LastLt(key K) (idx int) {
	return st.FirstGe(key) - 1
}

// LastLeはkey以下の最大要素のindexを返す。
// 存在しなければ-1を返す。要素はAt(LastLe(key))で取得。
func (st *Splaymap[K, V]) LastLe(key K) (idx int) {
	return st.FirstGt(key) - 1
}

// InOrderは昇順でキーと値のペアを列挙して返す
func (st *Splaymap[K, V]) InOrder() []Entry[K, V] {
	if st.root == nil {
		return nil
	}
	n := st.root.size
	ret := make([]Entry[K, V], 0, n)
	deq := NewDeque[*splaynode[K, V]]()

	cur := st.root
	for {
		if cur != nil {
			deq.PushBack(cur)
			cur = cur.l
			continue
		}
		if deq.Size() == 0 {
			break
		}
		cur = deq.PopBack()
		ret = append(ret, Entry[K, V]{cur.key, cur.value})
		cur = cur.r
	}
	return ret
}
