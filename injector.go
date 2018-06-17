package refinject

import "reflect"

type objEntry struct {
	l  label
	rv reflect.Value
}

// injector is a context of injection.
type injector struct {
	cat  *Catalog
	omap map[reflect.Type][]*objEntry
}

func newInjector(c *Catalog) *injector {
	return &injector{
		cat:  c,
		omap: make(map[reflect.Type][]*objEntry),
	}
}

func (j *injector) inject(rv reflect.Value) error {
	rv = reflect.Indirect(rv)
	typ := rv.Type()
	if typ.Kind() != reflect.Struct {
		return nil
	}
	num := typ.NumField()
	for i := 0; i < num; i++ {
		f := typ.Field(i)
		ityp, l, ok, err := needInject(f)
		if err != nil {
			return err
		}
		if !ok {
			// check embedded field.
			if fv, ok := isEmbedded(rv, f); ok {
				err := j.inject(fv)
				if err != nil {
					return err
				}
			}
			continue
		}
		if !rv.Field(i).CanSet() {
			return &CantSetError{rv: rv, i: i}
		}
		fv, err := j.materialize(ityp, l)
		if err != nil {
			// FIXME: include context of the error.
			return err
		}
		rv.Field(i).Set(fv)
	}
	return nil
}

func (j *injector) materialize(ityp reflect.Type, l label) (reflect.Value, error) {
	if ityp.Kind() != reflect.Interface {
		panic("type is not interface")
	}

	// find a type which implement ityp interface.
	typ, labels, err := j.cat.find(ityp, l)
	if err != nil {
		return reflect.Value{}, err
	}

	// create or re-use an object for typ.
	if rv, ok := j.cacheGet(typ, labels); ok {
		return rv, nil
	}
	rv := reflect.New(typ)
	j.cachePut(typ, labels, rv)

	// inject dependencies.
	err = j.inject(rv)
	if err != nil {
		return reflect.Value{}, err
	}
	return rv, nil
}

func (j *injector) cacheGet(ityp reflect.Type, l label) (reflect.Value, bool) {
	entries, ok := j.omap[ityp]
	if !ok {
		return reflect.Value{}, false
	}
	for _, e := range entries {
		if l.isSubset(e.l) {
			return e.rv, true
		}
	}
	return reflect.Value{}, false
}

func (j *injector) cachePut(ityp reflect.Type, l label, rv reflect.Value) {
	j.omap[ityp] = append(j.omap[ityp], &objEntry{l: l, rv: rv})
}
