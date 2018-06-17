package refinject

import (
	"reflect"
	"strings"
)

const (
	// Tag is used to mark fields as to be injected.
	Tag = "refinject"
)

func getType(v interface{}) (reflect.Type, error) {
	typ := reflect.TypeOf(v)
	for typ != nil && typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == nil {
		return nil, &NilTypeError{v: v}
	}
	return typ, nil
}

func getInterface(v interface{}) (reflect.Type, error) {
	typ, err := getType(v)
	if err != nil {
		return nil, err
	}
	if typ.Kind() != reflect.Interface {
		return nil, &NotInterfaceError{v: v}
	}
	return typ, nil
}

func needInject(f reflect.StructField) (reflect.Type, label, bool, error) {
	v, ok := f.Tag.Lookup(Tag)
	if !ok {
		return nil, nil, false, nil
	}
	if f.Type.Kind() != reflect.Interface {
		return nil, nil, false, &FieldNotInterfaceError{f: f}
	}
	return f.Type, newLabel(strings.Split(v, " ")), true, nil
}

func isEmbedded(rv reflect.Value, f reflect.StructField) (reflect.Value, bool) {
	if !f.Anonymous {
		return reflect.Value{}, false
	}
	fv := rv.FieldByIndex(f.Index)
	return fv, true
}

// DefaultCatalog is default injection catalog.
var DefaultCatalog = &Catalog{}

// Register registers a type to default catalog.
func Register(v interface{}, labels ...string) error {
	return DefaultCatalog.Register(v, labels...)
}

// Inject injects/fills fields which require to be injected from default
// catalog.
func Inject(v interface{}) error {
	return DefaultCatalog.Inject(v)
}

// Materialize creates an object which have implements an interface using
// default catalog.
func Materialize(v interface{}, labels ...string) (interface{}, error) {
	return DefaultCatalog.Materialize(v, labels...)
}
