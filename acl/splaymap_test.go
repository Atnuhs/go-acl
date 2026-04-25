package acl

import (
	"testing"
)

func TestSplayTree_BasicOperations(t *testing.T) {
	// 文字列キー、整数値のスプレー木
	tree := NewSplaymap[string, int]()

	// 挿入テスト
	tree.Insert("apple", 1)
	tree.Insert("banana", 2)
	tree.Insert("cherry", 3)

	if tree.Size() != 3 {
		t.Errorf("Expected size 3, got %d", tree.Size())
	}

	// 検索テスト
	if value, found := tree.Has("banana"); !found || value != 2 {
		t.Errorf("Expected to find banana with value 2, got value=%d, found=%v", value, found)
	}

	if _, found := tree.Has("orange"); found {
		t.Errorf("Expected not to find orange")
	}

	// 削除テスト
	if !tree.Delete("banana") {
		t.Errorf("Expected to delete banana successfully")
	}

	if tree.Size() != 2 {
		t.Errorf("Expected size 2 after deletion, got %d", tree.Size())
	}

	if _, found := tree.Has("banana"); found {
		t.Errorf("Expected banana to be deleted")
	}
}

func TestSplayTree_IntegerKeys(t *testing.T) {
	// 整数キー、文字列値のスプレー木
	tree := NewSplaymap[int, string]()

	// 複数の値を挿入
	values := map[int]string{
		5: "five",
		2: "two",
		8: "eight",
		1: "one",
		7: "seven",
	}

	for k, v := range values {
		tree.Insert(k, v)
	}

	// 中順巡回でソートされた順序を確認
	inOrder := tree.InOrder()
	expectedKeys := []int{1, 2, 5, 7, 8}

	if len(inOrder) != len(expectedKeys) {
		t.Errorf("Expected %d elements in in-order traversal, got %d", len(expectedKeys), len(inOrder))
	}

	for i, expected := range expectedKeys {
		if i >= len(inOrder) {
			t.Errorf("Missing element at index %d", i)
			continue
		}
		if inOrder[i].K != expected {
			t.Errorf("Expected key %d at index %d, got %d", expected, i, inOrder[i].K)
		}
		if inOrder[i].V != values[expected] {
			t.Errorf("Expected value %s for key %d, got %s", values[expected], expected, inOrder[i].V)
		}
	}
}

func TestSplayTree_UpdateValue(t *testing.T) {
	tree := NewSplaymap[string, int]()

	// 最初の値を挿入
	tree.Insert("key", 100)

	if value, found := tree.Has("key"); !found || value != 100 {
		t.Errorf("Expected to find key with value 100, got value=%d, found=%v", value, found)
	}

	// 同じキーで値を更新
	tree.Insert("key", 200)

	if tree.Size() != 1 {
		t.Errorf("Expected size 1 after update, got %d", tree.Size())
	}

	if value, found := tree.Has("key"); !found || value != 200 {
		t.Errorf("Expected to find key with updated value 200, got value=%d, found=%v", value, found)
	}
}

func TestSplayTree_EmptyOperations(t *testing.T) {
	tree := NewSplaymap[int, string]()

	if !tree.IsEmpty() {
		t.Errorf("Expected new tree to be empty")
	}

	if tree.Size() != 0 {
		t.Errorf("Expected size 0 for empty tree, got %d", tree.Size())
	}

	if _, found := tree.Has(1); found {
		t.Errorf("Expected not to find anything in empty tree")
	}

	if tree.Delete(1) {
		t.Errorf("Expected delete to return false for empty tree")
	}

	inOrder := tree.InOrder()
	if len(inOrder) != 0 {
		t.Errorf("Expected empty in-order traversal, got %d elements", len(inOrder))
	}
}

// パフォーマンステスト用のベンチマーク
func BenchmarkSplayTree_Insert(b *testing.B) {
	tree := NewSplaymap[int, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(i, i*2)
	}
}

func BenchmarkSplayTree_Find(b *testing.B) {
	tree := NewSplaymap[int, int]()

	// 事前に1000個の要素を挿入
	for i := 0; i < 1000; i++ {
		tree.Insert(i, i*2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Has(i % 1000)
	}
}
