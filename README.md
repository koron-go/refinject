# koron-go/refinject

[![GoDoc](https://godoc.org/github.com/koron-go/refinject?status.svg)](https://godoc.org/github.com/koron-go/refinject)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/refinject/master.svg)](https://circleci.com/gh/koron-go/refinject/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/refinject)](https://goreportcard.com/report/github.com/koron-go/refinject)

Experimental injection container by reflect.

## Getting started

This `refinject` package injects registered `struct`s to interfaces which
tagged as `refinject:"..."`.

### How `Inject` works

At first, there is a service provider which provides `Foo` method.

```go
package provider

type FooService struct {}

func (*Service) Foo() {}
```

And let's see how to inject.

```go
package main

import "github.com/koron-go/refinject"

type Fooer interface {
    Foo()
}

type FooUser struct {
    MyFoo Fooer `refinject:""`
}

func main() {
    refinject.Register(&provider.FooService{})
    
    var u FooUser
    err := refinject.Inject(&u)
    if err != nil {
        panic(err)
    }
    // work with u.MyFoo.Foo()
}
```

`Inject()` injects registered objects to each interface fields which tagged as
`refinject:""`. Injected objects will fulfill each interfaces.  If multiple
objects matches for an interface, injection will be failed.

The instance which passed to `Register()` is not used directly for injection.
New instance of same type will be injected.

### How `Materialize` works

Let's see how to use `Materialize`.

```go
package middle

type Fooer interface {
    Foo()
}

type BarService struct {
    MyFoo Fooer `refinject:""`
}

func (bar *BarService) Bar() {
    bar.MyFoo.Foo()
}
```

```go
package main

import "github.com/koron-go/refinject"

type Barer interface {
    Bar()
}

func main() {
    // register all components.
    refinject.Register(&provider.FooService{})
    refinject.Register(&middle.BarService{})

    var iv Barer
    _, err := refinject.Materialize(&iv)
    if err != nil {
        panic(err)
    }
    // work with iv.Bar()
}
```
