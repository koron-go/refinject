package refinject

import (
	"fmt"
	"reflect"
)

type typeEntry struct {
	typ reflect.Type
	ls  labelSet
}

func (e *typeEntry) String() string {
	return fmt.Sprintf("typeEntry:{%s labels:%+v}", e.typ, e.ls)
}

// Registry is a registry of components which provide interfaces to inject.
type Registry struct {
	typeMap map[reflect.Type]int
	entries []*typeEntry
}

// Register registers a type.
// The passed instance is not used, but only used its type,
// a new instance will be created when materialize.
func (reg *Registry) Register(v interface{}, labels ...string) error {
	typ, err := getType(v)
	if err != nil {
		return err
	}

	if reg.typeMap == nil {
		reg.typeMap = make(map[reflect.Type]int)
	}
	if n, ok := reg.typeMap[typ]; ok {
		old := reg.entries[n]
		return errorFunc(func() string {
			return fmt.Sprintf("registered already: %s old-label:%+v", typ, old.ls)
		})
	}
	reg.typeMap[typ] = len(reg.entries)
	reg.entries = append(reg.entries, &typeEntry{typ: typ, ls: newLabelSet(labels)})

	return nil
}

// Inject injects/fills fields which require to be injected by the component.
func (reg *Registry) Inject(v interface{}) error {
	return reg.inject(reflect.ValueOf(v))
}

func (reg *Registry) inject(rv reflect.Value) error {
	return newInjector(reg).inject(rv)
}

// Materialize materializes an object which have an interface ("v") with
// filling dependent interfaces.
func (reg *Registry) Materialize(v interface{}, labels ...string) (interface{}, error) {
	ityp, err := getInterface(v)
	if err != nil {
		return nil, err
	}
	rv, err := newInjector(reg).materialize(ityp, newLabelSet(labels))
	if err != nil {
		return nil, err
	}
	if !rv.CanInterface() {
		return nil, errorFunc(func() string {
			return fmt.Sprintf("won't be materialized: %s", rv)
		})
	}

	// FIXME: check to guard from panic
	iv := reflect.ValueOf(v).Elem()
	iv.Set(rv)

	return rv.Interface(), nil
}

// find finds a type which implements an interface (ityp) and match with
// labels.
func (reg *Registry) find(ityp reflect.Type, l labelSet) (reflect.Type, labelSet, error) {
	found := make([]*typeEntry, 0, 4)
	for _, e := range reg.entries {
		if !l.isSubset(e.ls) {
			continue
		}
		if reflect.PtrTo(e.typ).Implements(ityp) {
			found = append(found, e)
		}
	}
	switch len(found) {
	case 1:
		e := found[0]
		return e.typ, e.ls, nil
	case 0:
		return nil, nil, errorFunc(func() string {
			return fmt.Sprintf("not found interface: %s labels=%+v", ityp, l)
		})
	default:
		return nil, nil, errorFunc(func() string {
			return fmt.Sprintf("found multiple objects for interface:%s labels:%+v : %+v", ityp, l, found)
		})
	}
}
