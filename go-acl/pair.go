package main

import "fmt"

// Pair は2つの値を持つ構造体
type Pair[A any, B any] struct {
	U A
	V B
}

// NewPair Pairを生成する
func NewPair[A any, B any](u A, v B) *Pair[A, B] {
	return &Pair[A, B]{u, v}
}

// String Pairの文字列を、空白区切りで返す
func (p *Pair[A, B]) String() string {
	return fmt.Sprintf("%v %v", p.U, p.V)
}
