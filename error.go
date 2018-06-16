package refinject

import (
	"fmt"
	"reflect"
)

// NilTypeError means failed to get type from interface{}.
type NilTypeError struct {
	v interface{}
}

func (err *NilTypeError) Error() string {
	return fmt.Sprintf("failed to get type from: %#v", err.v)
}

// DupilcateTypeError means failed to register a type to a catalog.
type DupilcateTypeError struct {
	typ reflect.Type
}

func (err *DupilcateTypeError) Error() string {
	return fmt.Sprintf("type registered already: %s", err.typ)
}

// NotInterfaceError means passed value is not interface.
type NotInterfaceError struct {
	v interface{}
}

func (err *NotInterfaceError) Error() string {
	return fmt.Sprintf("not interface: %+v", err.v)
}

// NotFoundError means no providers for an interface.
type NotFoundError struct {
	ityp reflect.Type
	l    label
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("not found in catalog: interface=%s labels=%+v", err.ityp, err.l)
}

// CantMaterializeError means can't get interface{} from materialized value.
type CantMaterializeError struct {
	rv reflect.Value
}

func (err *CantMaterializeError) Error() string {
	return fmt.Sprintf("won't be materialized: %s", err.rv)
}

// CantSetError means non-public field is marked as to be injected.
type CantSetError struct {
	rv reflect.Value
	i  int
}

func (err *CantSetError) Error() string {
	typ := err.rv.Type()
	f := typ.Field(err.i)
	return fmt.Sprintf("won't be set: type=%s field=%s", typ, f.Name)
}

// FieldNotInterfaceError raises when not interface struct field is marked as
// to be injected.
type FieldNotInterfaceError struct {
	f reflect.StructField
}

func (err *FieldNotInterfaceError) Error() string {
	return fmt.Sprintf("non interface field: %s", err.f.Name)
}
