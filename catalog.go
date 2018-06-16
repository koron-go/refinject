package refinject

import "reflect"

// Catalog is a types catalog for injection.
type Catalog struct {
	tmap map[reflect.Type]label
}

// Register registers a type.
func (c *Catalog) Register(v interface{}, labels ...string) error {
	typ, err := getType(v)
	if err != nil {
		return err
	}
	if c.tmap == nil {
		c.tmap = make(map[reflect.Type]label)
	}
	if _, ok := c.tmap[typ]; ok {
		return &DupilcateTypeError{typ: typ}
	}
	c.tmap[typ] = newLabel(labels)
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
		return nil, &CantMaterializeError{rv: rv}
	}
	return rv.Interface(), nil
}

// find finds a type which implements an interface (ityp) and match with
// labels.
func (c *Catalog) find(ityp reflect.Type, l label) (reflect.Type, label, error) {
	for k, v := range c.tmap {
		if k.Implements(ityp) && l.isSubset(v) {
			return k, v, nil
		}
	}
	return nil, nil, &NotFoundError{ityp: ityp, l: l}
}
