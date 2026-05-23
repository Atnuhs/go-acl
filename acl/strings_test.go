package acl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestManacher(t *testing.T) {
	tests := map[string]struct {
		s    string
		want []int
	}{
		"odd":      {s: "ababa", want: []int{1, 2, 3, 2, 1}},
		"single":   {s: "a", want: []int{1}},
		"all_same": {s: "aaaa", want: []int{1, 2, 2, 1}},
		// 偶数長の回文を検知できるようにしていないもの
		"even": {s: "abba", want: []int{1, 1, 1, 1}},
		// 偶数長の回文を検知できるように入力文字を$で挟んで細工したもの
		"even2": {s: "$a$b$b$a$", want: []int{1, 2, 1, 2, 5, 2, 1, 2, 1}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Manacher(tc.s)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Manacher(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}

func TestKmpPrefix(t *testing.T) {
	tests := map[string]struct {
		s    string
		want []int
	}{
		"normal": {"aabaaabaab", []int{0, 1, 0, 1, 2, 2, 3, 4, 5, 3}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := kmpPrefix(tc.s)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("kmpPrefix(%q) mismatch (-want +got)\n%s", tc.s, diff)
			}
		})
	}
}
