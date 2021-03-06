package refinject

import (
	"strings"
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

func (bar *BarService) Bar() {
	bar.MyFoo.Foo()
}

type Barer interface {
	Bar()
}

type BazService struct {
	MyBar Barer `refinject:""`
}

type QuxService struct {
	MyFoo Fooer `refinject:""`
	MyBar Barer `refinject:""`
}

func TestInjectSimple(t *testing.T) {
	reg := &Registry{}
	reg.Register(&FooService{})

	bar := &BarService{}
	err := reg.Inject(bar)
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
	reg := &Registry{}
	reg.Register(&FooService{})

	var iv Fooer
	v, err := reg.Materialize(&iv)
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

func TestInjectHierarchy(t *testing.T) {
	reg := &Registry{}
	reg.Register(&FooService{})
	reg.Register(&BarService{})

	baz := &BazService{}
	err := reg.Inject(baz)
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

func TestInjectCached(t *testing.T) {
	reg := &Registry{}
	reg.Register(&FooService{})
	reg.Register(&BarService{})

	qux := &QuxService{}
	err := reg.Inject(qux)
	if err != nil {
		t.Fatalf("inject failed: %s", err)
	}

	pfoo, ok := qux.MyFoo.(*FooService)
	if !ok || pfoo == nil {
		t.Fatalf("failed to inject MyFoo: %+v", pfoo)
	}
	pbar, ok := qux.MyBar.(*BarService)
	if !ok || pbar == nil {
		t.Fatalf("failed to inject MyBar: %+v", pbar)
	}
	pfoo2, ok := pbar.MyFoo.(*FooService)
	if !ok || pfoo2 == nil {
		t.Fatalf("failed to inject MyBar.MyFoo: %+v", pfoo2)
	}
	if pfoo != pfoo2 {
		t.Fatalf("mismatch MyFoo and MyBar.MyFoo: %+v, %+v", pfoo, pfoo2)
	}
}

type Quuxer1 interface {
	Quux1()
}

type Quuxer2 interface {
	Quux2()
}

type QuuxService1 struct {
	MyQuux2 Quuxer2 `refinject:""`

	initCount int
}

func (*QuuxService1) Quux1() {}

func (q1 *QuuxService1) InitiateComponent() error {
	q1.initCount++
	return nil
}

type QuuxService2 struct {
	MyQuux1 Quuxer1 `refinject:""`
}

func (*QuuxService2) Quux2() {}

func TestMaterializeRecursive(t *testing.T) {
	reg := &Registry{}
	reg.Register(&QuuxService1{})
	reg.Register(&QuuxService2{})

	var iv Quuxer1
	_, err := reg.Materialize(&iv)
	if err != nil {
		t.Fatalf("materialize failed: %s", err)
	}

	pq1, ok := iv.(*QuuxService1)
	if !ok || pq1 == nil {
		t.Fatalf("failed to materialize Quuxer1: %+v", pq1)
	}
	pq2, ok := pq1.MyQuux2.(*QuuxService2)
	if !ok || pq2 == nil {
		t.Fatalf("failed to inject MyQuux2: %+v", pq2)
	}
	pq1b, ok := pq2.MyQuux1.(*QuuxService1)
	if !ok || pq2 == nil {
		t.Fatalf("failed to inject MyQuux1: %+v", pq1b)
	}
	if pq1b != pq1 {
		t.Fatalf("faield to re-use Quuxer1: %p %p", pq1b, pq1)
	}

	t.Run("initiator", func(t *testing.T) {
		if pq1.initCount != 1 {
			t.Fatalf("unexpected call of InitiateComponent: %d", pq1.initCount)
		}
	})
}

type CorgeService struct {
	BarService
}

func TestInjectEmbedded(t *testing.T) {
	reg := &Registry{}
	reg.Register(&CorgeService{})
	reg.Register(&FooService{})

	baz := &BazService{}
	err := reg.Inject(baz)
	if err != nil {
		t.Fatalf("inject failed: %s", err)
	}

	p, ok := baz.MyBar.(*CorgeService)
	if !ok || p == nil {
		t.Fatalf("failed to inject MyFoo: %+v", p)
	}
	// assert not panic
	p.Bar()
}

func TestMaterializeEmbedded(t *testing.T) {
	reg := &Registry{}
	reg.Register(&CorgeService{})
	reg.Register(&FooService{})

	var iv Barer
	v, err := reg.Materialize(&iv)
	if err != nil {
		t.Fatalf("materialize failed: %s", err)
	}
	p, ok := v.(*CorgeService)
	if !ok || p == nil {
		t.Fatalf("failed to materialize Fooer: %+v", v)
	}
	// assert not panic
	iv.Bar()
}

func TestFoundMultipleObjects(t *testing.T) {
	reg := &Registry{}
	reg.Register(&BarService{})
	reg.Register(&CorgeService{})

	var iv Barer
	v, err := reg.Materialize(&iv)
	if err == nil {
		t.Fatalf("materialize should be failed, by multiple found: %+v", v)
	}
	if !strings.HasPrefix(err.Error(), "found multiple objects for ") {
		t.Errorf("unexpected error message: %s", err)
	}
}

type FooService2 struct {
	initCount int
}

func (*FooService2) Foo() {}

func (foo *FooService2) InitiateComponent() error {
	foo.initCount++
	return nil
}

func TestInitiator(t *testing.T) {
	reg := &Registry{}
	reg.Register(&FooService2{})

	var iv Fooer
	_, err := reg.Materialize(&iv)
	if err != nil {
		t.Fatalf("materialize failed: %s", err)
	}
	p, ok := iv.(*FooService2)
	if !ok || p == nil {
		t.Fatalf("unexpected return value: %+v", p)
	}
	if p.initCount != 1 {
		t.Fatalf("unexpected call of InitiateComponent: %d", p.initCount)
	}
}
