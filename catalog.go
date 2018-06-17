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

// Catalog is a types catalog for injection.
type Catalog struct {
	typeMap  map[reflect.Type]int
	entries []*typeEntry
}

// Register registers a type.
func (c *Catalog) Register(v interface{}, labels ...string) error {
	typ, err := getType(v)
	if err != nil {
		return err
	}

	if c.typeMap == nil {
		c.typeMap = make(map[reflect.Type]int)
	}
	if n, ok := c.typeMap[typ]; ok {
		old := c.entries[n]
		return errorFunc(func() string {
			return fmt.Sprintf("registered already: %s old-label:%+v", typ, old.ls)
		})
	}
	c.typeMap[typ] = len(c.entries)
	c.entries = append(c.entries, &typeEntry{typ: typ, ls: newLabel(labels)})

	return nil
}

// Inject injects/fills dependent interfaces of v from the catalog.
func (c *Catalog) Inject(v interface{}) error {
	return c.inject(reflect.ValueOf(v))
}

func (c *Catalog) inject(rv reflect.Value) error {
	return newInjector(c).inject(rv)
}

// Materialize materializes an object which have an interface ("v") with
// filling dependent interfaces.
func (c *Catalog) Materialize(v interface{}, labels ...string) (interface{}, error) {
	ityp, err := getInterface(v)
	if err != nil {
		return nil, err
	}
	rv, err := newInjector(c).materialize(ityp, newLabel(labels))
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
func (c *Catalog) find(ityp reflect.Type, l labelSet) (reflect.Type, labelSet, error) {
	found := make([]*typeEntry, 0, 4)
	for _, e := range c.entries {
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
