package refinject

import "testing"

func TestLabel(t *testing.T) {
	isSubset := func(a []string, b []string) {
		la := newLabel(a)
		lb := newLabel(b)
		if !la.isSubset(lb) {
			t.Errorf("isSubset failed:\n\ta=%+v\n\tb=%+v\n\tla=%+v\n\tlb=%+v",
				a, b, la, lb)
		}
	}

	// subset: empty label
	isSubset(nil, nil)
	isSubset([]string{}, nil)
	isSubset(nil, []string{})
	isSubset([]string{}, []string{})

	// subset: basic
	p1 := []string{"aaa", "bbb", "ccc"}
	isSubset(nil, p1)
	isSubset([]string{"aaa"}, p1)
	isSubset([]string{"bbb"}, p1)
	isSubset([]string{"ccc"}, p1)
	isSubset([]string{"aaa", "bbb"}, p1)
	isSubset([]string{"aaa", "ccc"}, p1)
	isSubset([]string{"bbb", "ccc"}, p1)
	isSubset([]string{"aaa", "bbb", "ccc"}, p1)

	// subset: non-sorted super set
	p2 := []string{"ccc", "bbb", "aaa"}
	isSubset(nil, p2)
	isSubset([]string{"aaa"}, p2)
	isSubset([]string{"bbb"}, p2)
	isSubset([]string{"ccc"}, p2)
	isSubset([]string{"aaa", "bbb"}, p2)
	isSubset([]string{"aaa", "ccc"}, p2)
	isSubset([]string{"bbb", "ccc"}, p2)
	isSubset([]string{"aaa", "bbb", "ccc"}, p2)

	// subset: non-sorted sub-set
	isSubset([]string{"bbb", "aaa"}, p1)
	isSubset([]string{"ccc", "aaa"}, p1)
	isSubset([]string{"ccc", "bbb"}, p1)
	isSubset([]string{"aaa", "ccc", "bbb"}, p1)
	isSubset([]string{"bbb", "aaa", "ccc"}, p1)
	isSubset([]string{"bbb", "ccc", "aaa"}, p1)
	isSubset([]string{"ccc", "aaa", "bbb"}, p1)
	isSubset([]string{"ccc", "bbb", "aaa"}, p1)

	// subset: non-sorted both
	isSubset([]string{"bbb", "aaa"}, p2)
	isSubset([]string{"ccc", "aaa"}, p2)
	isSubset([]string{"ccc", "bbb"}, p2)
	isSubset([]string{"aaa", "ccc", "bbb"}, p2)
	isSubset([]string{"bbb", "aaa", "ccc"}, p2)
	isSubset([]string{"bbb", "ccc", "aaa"}, p2)
	isSubset([]string{"ccc", "aaa", "bbb"}, p2)
	isSubset([]string{"ccc", "bbb", "aaa"}, p2)

	// subset: ignore empty label
	isSubset([]string{""}, nil)
	isSubset([]string{""}, []string{})
	isSubset([]string{""}, p1)
	isSubset([]string{""}, p2)

	nonSubset := func(a []string, b []string) {
		la := newLabel(a)
		lb := newLabel(b)
		if la.isSubset(lb) {
			t.Errorf("isSubset unexpectedly success:\n\ta=%+v\n\tb=%+v\n\tla=%+v\n\tlb=%+v",
				a, b, la, lb)
		}
	}

	// non-subset: empty label
	nonSubset([]string{"xxx"}, nil)
	nonSubset([]string{"xxx", "yyy"}, nil)
	nonSubset([]string{"xxx", "yyy", "zzz"}, nil)

	// non-subset: basic
	nonSubset([]string{"xxx"}, p1)
	nonSubset([]string{"yyy"}, p1)
	nonSubset([]string{"zzz"}, p1)
	nonSubset([]string{"xxx", "yyy"}, p1)
	nonSubset([]string{"xxx", "zzz"}, p1)
	nonSubset([]string{"yyy", "zzz"}, p1)
	nonSubset([]string{"xxx", "yyy", "zzz"}, p1)
}
