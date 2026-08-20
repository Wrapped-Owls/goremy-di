---
title: "Decorating Binds"
weight: 7
menu:
  main:
    parent: advanced
    params:
      icon: "🎁"
---

`Decorate` wraps a bind so a decorator runs over the value it produces, without touching the registration the bind came
from. It is the classic decorator: auditing, tracing, metrics or retry added around a dependency that knows nothing
about them.

Binds nobody decorates pay nothing.

---

## Wrapping a bind you do not own

A module exposes the bind; the application decides whether to wrap it:

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

// Owned by the module
func greeterBind() remy.Bind[Greeter] {
	return remy.Singleton(func(remy.DependencyRetriever) (Greeter, error) {
		return EnglishGreeter{}, nil
	})
}

// Owned by the application
func main() {
	injector := remy.NewInjector()
	remy.Register(injector, remy.Decorate[Greeter](greeterBind(), shout))
}

func shout(_ remy.DependencyRetriever, inner Greeter) (Greeter, error) {
	return shoutingGreeter{inner: inner}, nil
}
```

The consumer keeps calling `remy.Get[Greeter]` and never learns a wrapper is in the way.

---

## Building the wrapper from other dependencies

A `Decorator[T]` receives the `DependencyRetriever`, so the wrapper can resolve its own collaborators from the
container instead of having them handed in:

```go
func auditGreeter(retriever remy.DependencyRetriever, inner Greeter) (Greeter, error) {
	audit, err := remy.Get[*AuditLog](retriever)
	if err != nil {
		return nil, err
	}

	return auditedGreeter{inner: inner, audit: audit}, nil
}

remy.RegisterInstance(injector, &AuditLog{})
remy.Register(injector, remy.Decorate[Greeter](greeterBind(), auditGreeter))
```

An error from the decorator, including a collaborator it could not resolve, reaches the caller of `Get` unchanged.

---

## Stacking

Decorators compose from the inside out, so the last one applied is the outermost:

```go
inner := remy.Decorate[Greeter](greeterBind(), audit)
remy.Register(injector, remy.Decorate[Greeter](inner, retry))
// retry wraps audit, audit wraps the greeter
```

---

## Lifecycle

The decoration is folded into the bind's own generator rather than layered on top of it, so it shares whatever
lifecycle the bind already had:

| Wrapped bind     | Times the decorator runs               |
|------------------|----------------------------------------|
| `Instance`       | Once, while it registers               |
| `Singleton`      | Once, while it registers               |
| `LazySingleton`  | Once, on the first retrieval           |
| `Factory`        | Every value it builds                  |

This matters when the decorator allocates a wrapper: a decorated singleton hands out the same wrapper every time,
exactly as an undecorated one hands out the same value.
