package refinject

import "sort"

type labelSet []string

var emptyLabel = labelSet{}

func newLabelSet(v []string) labelSet {
	if len(v) == 0 {
		return emptyLabel
	}
	w := make([]string, 0, len(v))
	m := make(map[string]struct{})
	for _, s := range v {
		if s == "" {
			continue
		}
		if _, ok := m[s]; ok {
			continue
		}
		m[s] = struct{}{}
		w = append(w, s)
	}
	sort.Strings(w)
	return labelSet(w)
}

func (l labelSet) isSubset(p labelSet) bool {
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
