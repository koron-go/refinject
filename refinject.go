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

func needInject(f reflect.StructField) (reflect.Type, labelSet, bool, error) {
	v, ok := f.Tag.Lookup(Tag)
	if !ok {
		return nil, nil, false, nil
	}
	if f.Type.Kind() != reflect.Interface {
		return nil, nil, false, errorFunc(func() string {
			return fmt.Sprintf("non interface field won't be injected: %s", f.Name)
		})
	}
	return f.Type, newLabelSet(strings.Split(v, " ")), true, nil
}

func isEmbedded(rv reflect.Value, f reflect.StructField) (reflect.Value, bool) {
	if !f.Anonymous {
		return reflect.Value{}, false
	}
	fv := rv.FieldByIndex(f.Index)
	return fv, true
}

// DefaultRegistry is default components registry.
var DefaultRegistry = &Registry{}

// Register registers a component's type to default registry.
// The passed instance is not used, but only used its type.
// A new instance will be created when materialize.
func Register(v interface{}, labels ...string) error {
	return DefaultRegistry.Register(v, labels...)
}

// Inject injects/fills fields which require to be injected by the component.
func Inject(v interface{}) error {
	return DefaultRegistry.Inject(v)
}

// Materialize creates a new component which fulfills required interfaces.
// A component will be injected by all dependencies.
func Materialize(v interface{}, labels ...string) (interface{}, error) {
	return DefaultRegistry.Materialize(v, labels...)
}

// Initiator is called when a component is created and all fields injected.
type Initiator interface {
	InitiateComponent() error
}

func initiateComponent(rv reflect.Value) error {
	if !rv.CanInterface() {
		return nil
	}
	p, ok := rv.Interface().(Initiator)
	if !ok {
		return nil
	}
	err := p.InitiateComponent()
	if err != nil {
		return errorFunc(func() string {
			return fmt.Sprintf("failed to initiate component %T: %s", p, err)
		})
	}
	return nil
}
