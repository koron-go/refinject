package refinject

import "sort"

type label []string

var emptyLabel = label{}

func newLabel(v []string) label {
	if len(v) == 0 {
		return emptyLabel
	}
	w := make([]string, len(v))
	copy(w, v)
	sort.Strings(w)
	return label(w)
}

func (l label) isSubset(p label) bool {
	if len(l) == 0 {
		return true
	}
	if len(p) == 0 {
		return false
	}
	i := 0
	for _, s := range l {
		for ; i < len(p); i++ {
			t := p[i]
			if s == t {
				break
			}
			if s < t {
				return false
			}
		}
		if i >= len(p) {
			return false
		}
		i++
	}
	return true
}
