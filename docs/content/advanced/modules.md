---
title: "Modules"
weight: 5
---

Modules let you group multiple registrations and apply them to an injector in a single call. They are a thin, optional
abstraction built on top of the existing API; you can keep using the regular Register* functions directly, or compose
them via modules for cleaner setup code.

Why use modules?

- Encapsulate registrations per feature or layer (e.g., http, repo, service)
- Build small modules and compose them in your app wiring
- Keep startup code tidy and declarative

Core concepts

- Module: an object with Register(Injector) that applies registrations.
- NewModule(registers ...ModuleRegister): creates a module from a list of registration functions.
- ModuleRegister: a function that receives an Injector and performs registrations.

Quick start

```go
package app

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

type UserSvc struct{ Greeting string }

func NewUserSvc(greeting string) (UserSvc, error) {
	return UserSvc{Greeting: greeting}, nil
}

var userModule = remy.NewModule(
	remy.WithInstance("hello"),
	remy.WithConstructor(remy.Factory[UserSvc], NewUserSvc),
)

func Wire(inj remy.Injector) error {
	return remy.RegisterModule(inj, userModule)
}
```

Option helpers (With*)

The With* helpers adapt the existing registration API into module-ready registration functions. Only error-returning
constructor adapters are provided.

- WithInstance[T any](value T, optTag ...string)
- WithBind[T any](bind remy.Bind[T], optTag ...string)
- WithFactory[T any](binder types.Binder[T], optTag ...string)
- WithSingleton[T any](binder types.Binder[T], optTag ...string)
- WithLazySingleton[T any](binder types.Binder[T], optTag ...string)
- WithConstructor[T any](Factory[T], func() (T, error), optTag ...string)
- WithConstructor1[T, A any](Factory[T], func(A) (T, error), optTag ...string)
- WithConstructor2[T, A, B any](Factory[T], func(A, B) (T, error), optTag ...string)
- WithConstructor3[T, A, B, C any](Factory[T], func(A, B, C) (T, error), optTag ...string)
- WithConstructor4[T, A, B, C, D any](Factory[T], func(A, B, C, D) (T, error), optTag ...string)

Composing modules

You can register multiple modules at once. They are applied in the order provided.

```go
base := remy.NewModule(
    remy.WithInstance("world"),
)

derived := remy.NewModule(
    remy.WithConstructor1(remy.Factory[int], func(s string) (int, error) {
        return len(s), nil
    }),
)

inj := remy.NewInjector()
if err := remy.RegisterModule(inj, base, derived); err != nil {
    panic(err)
}

_ = remy.MustGet[int](inj) // 5
```

Error handling and panics

- The underlying Register* functions may panic on invalid actions (e.g., duplicate registration when override is
  disabled).
- RegisterModule recovers those panics and returns them as error, stopping further module options from running.
- Any successful registrations performed before a panic remain applied.

Interoperability

Modules do not replace the standard API. They simply let you stage registrations and apply them later. You can mix
direct registrations and modules freely:

```go
inj := remy.NewInjector()

// Direct registration
remy.RegisterInstance(inj, 42)

// Module-based registrations
svcMod := remy.NewModule(
    remy.WithFactory(remy.LazySingleton(func(_ remy.DependencyRetriever) (string, error) {
        return "data", nil
    })),
)
_ = remy.RegisterModule(inj, svcMod)
```

Naming notes

- The preferred helpers are the With* functions, following Go’s functional options style.

Tips for organizing modules

- Group by feature: user.Module, billing.Module, http.Module
- Provide small, reusable modules (e.g., logging, config, database)
- Compose them in your application’s main wiring function

Related API

- NewModule(registers ...ModuleRegister) Module
- RegisterModule(inj Injector, modules ...Module) error
- RegisterModuleFunc(inj Injector, regs ...func(Injector)) error
- All standard Register*, Override*, and Factory/Singleton helpers remain available
