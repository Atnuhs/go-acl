package acl

import (
	"cmp"
	"fmt"
)

type splitMode int

const (
	SplitLT_GE splitMode = iota // L: {x | x.key < key }, R: {x | x.key >= key}
	SplitLE_GT                  // L: {x | x.key <= key }, R: {x | x.key > key}
)

// splaynode はスプレー木のノード
type splaynode[K cmp.Ordered, V any] struct {
	l, r, p *splaynode[K, V]
	key     K
	value   V
	size    int
}

// ---- ノードのメソッド ----

func (n *splaynode[K, V]) update() {
	if n == nil {
		return
	}
	n.size = 1
	if n.l != nil {
		n.size += n.l.size
	}
	if n.r != nil {
		n.size += n.r.size
	}
}

func (n *splaynode[K, V]) rotate() {
	p := n.p
	if p == nil {
		return
	}

	pp := p.p
	if pp != nil {
		if pp.l == p {
			pp.l = n
		} else {
			pp.r = n
		}
	}
	n.p = pp

	if p.l == n {
		p.l = n.r
		if n.r != nil {
			n.r.p = p
		}
		n.r = p
	} else {
		p.r = n.l
		if n.l != nil {
			n.l.p = p
		}
		n.l = p
	}
	p.p = n

	p.update()
	n.update()
}

func (n *splaynode[K, V]) splay() {
	// スプレー操作の実装（既存のSplayNodeと同様）
	for n.p != nil {
		if n.p.p == nil {
			// Zig step
			n.rotate()
		} else if (n.p.l == n) == (n.p.p.l == n.p) {
			// Zig-zig step
			n.p.rotate()
			n.rotate()
		} else {
			// Zig-zag step
			n.rotate()
			n.rotate()
		}
	}
}

func (n *splaynode[K, V]) splayMax() *splaynode[K, V] {
	if n == nil {
		return nil
	}
	cur := n
	for cur.r != nil {
		cur = cur.r
	}
	cur.splay()
	return cur
}

func (n *splaynode[K, V]) splayMin() *splaynode[K, V] {
	if n == nil {
		return nil
	}
	cur := n
	for cur.l != nil {
		cur = cur.l
	}
	cur.splay()
	return cur
}

func (n *splaynode[K, V]) has(key K) (node *splaynode[K, V], found bool) {
	if n == nil {
		return nil, false
	}

	node = n
	var last *splaynode[K, V]
	for node != nil {
		last = node
		if key == node.key {
			node.splay()
			return node, true
		} else if key < node.key {
			node = node.l
		} else {
			node = node.r
		}
	}
	// not found
	last.splay()
	return last, false
}

func (n *splaynode[K, V]) kth(k int) *splaynode[K, V] {
	if n == nil {
		return nil
	}

	if k < 0 || k >= n.size {
		return nil
	}

	cur := n
	for {
		curK := 0
		if cur.l != nil {
			curK = cur.l.size
		}
		switch {
		case k < curK:
			cur = cur.l
		case k == curK:
			cur.splay()
			return cur
		default:
			k -= curK + 1
			cur = cur.r
		}
	}
}
func (p *splaynode[K, V]) cutRight() (L, R *splaynode[K, V]) {
	if p == nil {
		return nil, nil
	}
	L, R = p, p.r
	if R != nil {
		R.p = nil
	}
	L.r = nil
	L.update()
	return
}

func (p *splaynode[K, V]) cutLeft() (L, R *splaynode[K, V]) {
	if p == nil {
		return nil, nil
	}
	L, R = p.l, p
	if L != nil {
		L.p = nil
	}
	R.l = nil
	R.update()
	return
}

// split
// splitMode=SplitLE_GT => (L {x | x.key <  key}, R: {x | x.key >= key})
// splitMode=SplitLT_GE => (L {x | x.key <= key}, R: {x | x.key >  key})
func split[K cmp.Ordered, V any](root *splaynode[K, V], key K, mode splitMode) (L, R *splaynode[K, V]) {
	if root == nil {
		return nil, nil
	}
	root, found := root.has(key)
	if found {
		switch mode {
		case SplitLE_GT:
			return root.cutRight()
		case SplitLT_GE:
			return root.cutLeft()
		default:
			panic(fmt.Sprintf("invalid splitMode: %+v", mode))
		}
	}

	// if not found
	if root.key < key {
		return root.cutRight()
	}
	return root.cutLeft()
}

func splitAt[V any](root *splaynode[int, V], i int) (L, R *splaynode[int, V]) {
	if root == nil {
		return
	}
	if i < 0 {
		return nil, root
	}
	if i >= root.size {
		return root, nil
	}
	root = root.kth(i)
	return root.cutLeft()
}

func merge[K cmp.Ordered, V any](L, R *splaynode[K, V]) *splaynode[K, V] {
	if L == nil {
		return R
	}
	if R == nil {
		return L
	}

	L = L.splayMax()
	L.r = R
	R.p = L
	L.update()
	return L
}
