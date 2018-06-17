package refinject

import (
	"testing"
)

type Fooer interface {
	Foo()
}

type FooService struct{}

func (*FooService) Foo() {}

var _ Fooer = (*FooService)(nil)

type BarService struct {
	MyFoo Fooer `refinject:""`
}

func (*BarService) Bar() {}

type Barer interface {
	Bar()
}

type BazService struct {
	MyBar Barer `refinject:""`
}

func TestInjectSimple(t *testing.T) {
	c := &Catalog{}
	c.Register(&FooService{})

	bar := &BarService{}
	err := c.Inject(bar)
	if err != nil {
		t.Fatalf("inject failed: %s", err)
	}
	p, ok := bar.MyFoo.(*FooService)
	if !ok {
		t.Fatalf("injected unexpected: %+v", p)
	}
	if p == nil {
		t.Fatalf("injected nil")
	}
}

func TestMaterializeSimple(t *testing.T) {
	c := &Catalog{}
	c.Register(&FooService{})

	var iv Fooer
	v, err := c.Materialize(&iv)
	if err != nil {
		t.Fatalf("materialize failed: %s", err)
	}
	p0, ok := v.(*FooService)
	if !ok || p0 == nil {
		t.Fatalf("unexpected return value: %+v", p0)
	}
	p1, ok := iv.(*FooService)
	if !ok || p1 == nil {
		t.Fatalf("unexpected out-arg: %+v", p1)
	}
	if p0 != p1 {
		t.Fatalf("not match return and out-arg values: %+v %+v", p0, p1)
	}
}

func TestInjectRecur(t *testing.T) {
	c := &Catalog{}
	c.Register(&FooService{})
	c.Register(&BarService{})

	baz := &BazService{}
	err := c.Inject(baz)
	if err != nil {
		t.Fatalf("inject failed: %s", err)
	}
	pbar, ok := baz.MyBar.(*BarService)
	if !ok || pbar == nil {
		t.Fatalf("failed injection Barer: %+v", pbar)
	}
	pfoo, ok := pbar.MyFoo.(*FooService)
	if !ok || pfoo == nil {
		t.Fatalf("failed injection Fooer: %+v", pfoo)
	}
}
