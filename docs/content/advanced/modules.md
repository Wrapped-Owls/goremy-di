---
title: "Modules"
weight: 5
menu:
  main:
    parent: advanced
    params:
      icon: "🧩"
---

Modules let you group multiple registrations and apply them to an injector in a single call. They are a thin, optional
abstraction built on top of the existing API; you can keep using the regular `Register*` functions directly, or compose
them via modules for cleaner setup code.

**Key Points:**

- 🧩 Group related registrations together
- 🔧 Encapsulate registrations per feature or layer
- 🎯 Keep startup code tidy and declarative
- 🔄 Compose small modules into larger applications

---

## Core Concepts

- **Module**: An object with `Register(Injector)` that applies registrations
- **NewModule**: Creates a module from a list of registration functions
- **ModuleRegister**: A function that receives an `Injector` and performs registrations

---

## Basic Example

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

---

## Option Helpers (With\*)

The `With*` helpers adapt the existing registration API into module-ready registration functions. Only error-returning
constructor adapters are provided.

### Available Helpers

- `WithInstance[T any](value T, optTag ...string)`
- `WithBind[T any](bind remy.Bind[T], optTag ...string)`
- `WithFactory[T any](binder types.Binder[T], optTag ...string)`
- `WithSingleton[T any](binder types.Binder[T], optTag ...string)`
- `WithLazySingleton[T any](binder types.Binder[T], optTag ...string)`
- `WithConstructor[T any](Factory[T], func() (T, error), optTag ...string)`
- `WithConstructor1[T, A any](Factory[T], func(A) (T, error), optTag ...string)`
- `WithConstructor2[T, A, B any](Factory[T], func(A, B) (T, error), optTag ...string)`
- `WithConstructor3[T, A, B, C any](Factory[T], func(A, B, C) (T, error), optTag ...string)`
- `WithConstructor4[T, A, B, C, D any](Factory[T], func(A, B, C, D) (T, error), optTag ...string)`

> **ℹ️ INFO:** The preferred helpers are the `With*` functions, following Go's functional options style.

---

## Composing Modules

You can register multiple modules at once. They are applied in the order provided.

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
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
}
```

---

## Error Handling

- The underlying `Register*` functions may panic on invalid actions (e.g., duplicate registration when override is
  disabled)
- `RegisterModule` recovers those panics and returns them as error, stopping further module options from running
- Any successful registrations performed before a panic remain applied

```go
package main

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	inj := remy.NewInjector(remy.Config{
		CanOverride: false, // Default
	})

	module := remy.NewModule(
		remy.WithInstance("first"),
		remy.WithInstance("second"), // This will cause an error
	)

	if err := remy.RegisterModule(inj, module); err != nil {
		// Error is returned, first registration remains applied
		log.Fatal(err)
	}
}
```

---

## Interoperability

Modules do not replace the standard API. They simply let you stage registrations and apply them later. You can mix
direct registrations and modules freely:

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
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

	// Both are available
	_ = remy.MustGet[int](inj)    // 42
	_ = remy.MustGet[string](inj) // "data"
}
```

---

## Best Practices

### Organizing Modules

- **Group by feature**: `user.Module`, `billing.Module`, `http.Module`
- **Provide small, reusable modules**: e.g., logging, config, database
- **Compose them in your application's main wiring function**

### Example: Feature-Based Organization

```go
// user/user_module.go
package user

import "github.com/wrapped-owls/goremy-di/remy"

var Module = remy.NewModule(
	remy.WithSingleton(remy.Factory[*UserService], NewUserService),
	remy.WithSingleton(remy.Factory[*UserRepository], NewUserRepository),
)
```

```go
// billing/billing_module.go
package billing

import "github.com/wrapped-owls/goremy-di/remy"

var Module = remy.NewModule(
	remy.WithSingleton(remy.Factory[*BillingService], NewBillingService),
)
```

```go
// main.go
package main

import (
	"app/billing"
	"app/user"
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	inj := remy.NewInjector()

	// Compose all modules
	if err := remy.RegisterModule(inj, user.Module, billing.Module); err != nil {
		panic(err)
	}

	// Application code...
}
```

---

## Related API

- `NewModule(registers ...ModuleRegister) Module` - Creates a new module from registration functions
- `RegisterModule(inj Injector, modules ...Module) error` - Registers one or more modules on an injector
- `RegisterModuleFunc(inj Injector, modules ...func(Injector)) error` - Registers modules using raw functions

---

## Module Interface

You can also implement the `Module` interface directly if you need more control:

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

type Module interface {
	Register(injector remy.Injector)
}

```

This allows you to create custom module implementations beyond the `NewModule` helper:

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

type CustomModule struct {
	// Your custom fields
}

func (m *CustomModule) Register(injector remy.Injector) {
	// Your custom registration logic
	remy.RegisterInstance(injector, "custom-value")
}

func main() {
	// Usage
	customMod := &CustomModule{}
	_ = remy.RegisterModule(inj, customMod)
}

```
