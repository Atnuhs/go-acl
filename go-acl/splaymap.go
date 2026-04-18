package main

import "cmp"

// Splaymap は汎用的なスプレー木の実装
type Splaymap[K cmp.Ordered, V any] struct {
	root *splaynode[K, V]
}

// NewSplaymap は新しいスプレー木を作成
func NewSplaymap[K cmp.Ordered, V any]() *Splaymap[K, V] {
	return &Splaymap[K, V]{}
}

// Size は要素数を返す
func (st *Splaymap[K, V]) Size() int {
	if st.root == nil {
		return 0
	}
	return st.root.size
}

// IsEmpty は空かどうかを返す
func (st *Splaymap[K, V]) IsEmpty() bool {
	return st.Size() == 0
}

// Has はキーで検索し、値と存在フラグを返す
func (st *Splaymap[K, V]) Has(key K) (value V, found bool) {
	if st.root == nil {
		return
	}

	st.root, found = st.root.has(key)
	if found {
		value = st.root.value
	}
	return
}

// Insert はキーと値のペアを挿入
func (st *Splaymap[K, V]) Insert(key K, value V) {
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

// Delete はキーで要素を削除
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

// Kthは[0, st.Size())の範囲内でk番目の要素を返す。
func (st *Splaymap[K, V]) Kth(k int) (kk K, vv V, ok bool) {
	if st.root == nil {
		return
	}

	if k < 0 || k >= st.Size() {
		return
	}

	st.root = st.root.kth(k)
	return st.root.key, st.root.value, true
}

func (st *Splaymap[K, V]) GtAt(key K) (idx int) {
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

func (st *Splaymap[K, V]) GeAt(key K) (idx int) {
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

func (st *Splaymap[K, V]) LtAt(key K) (idx int) {
	idx = st.GeAt(key) - 1
	return
}
func (st *Splaymap[K, V]) LeAt(key K) (idx int) {
	idx = st.GtAt(key) - 1
	return
}

func (st *Splaymap[K, V]) Le(key K) (k K, v V, ok bool) {
	return st.Kth(st.LeAt(key))
}

func (st *Splaymap[K, V]) Lt(key K) (k K, v V, ok bool) {
	return st.Kth(st.LtAt(key))
}

func (st *Splaymap[K, V]) Ge(key K) (k K, v V, ok bool) {
	return st.Kth(st.GeAt(key))
}

func (st *Splaymap[K, V]) Gt(key K) (k K, v V, ok bool) {
	return st.Kth(st.GtAt(key))
}

// InOrder は中順巡回でキーと値のペアを返す
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
