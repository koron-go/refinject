package refinject

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// Tag is used to mark fields as to be injected.
	Tag = "refinject"
)

type errorFunc func() string

func (f errorFunc) Error() string {
	return f()
}

func getType(v interface{}) (reflect.Type, error) {
	typ := reflect.TypeOf(v)
	for typ != nil && typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == nil {
		return nil, errorFunc(func() string {
			return fmt.Sprintf("failed to determine type of: %+v", v)
		})
	}
	return typ, nil
}

func getInterface(v interface{}) (reflect.Type, error) {
	typ, err := getType(v)
	if err != nil {
		return nil, err
	}
	if typ.Kind() != reflect.Interface {
		return nil, errorFunc(func() string {
			return fmt.Sprintf("failed to determine interface: %+v", v)
		})
	}
	return typ, nil
}

func needInject(f reflect.StructField) (reflect.Type, label, bool, error) {
	v, ok := f.Tag.Lookup(Tag)
	if !ok {
		return nil, nil, false, nil
	}
	if f.Type.Kind() != reflect.Interface {
		return nil, nil, false, errorFunc(func() string {
			return fmt.Sprintf("non interface field won't be injected: %s", f.Name)
		})
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
