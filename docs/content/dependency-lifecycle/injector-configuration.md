---
title: "Injector Configuration"
weight: 4
menu:
  main:
    parent: dependency-lifecycle
    params:
      icon: "⚙️"
---

Remy DI provides flexible injector configuration options to customize how dependencies are registered, resolved, and
managed. You can create multiple injector instances with different configurations to suit your application's needs.

## Creating an Injector

The `NewInjector` function creates a new injector instance. You can optionally pass a `Config` struct to customize its
behavior:

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	// Create injector with default configuration
	injector := remy.NewInjector()

	// Create injector with custom configuration
	customInjector := remy.NewInjector(remy.Config{
		CanOverride:      true,
		DuckTypeElements: false,
	})
}
```

## Configuration Options

### CanOverride

**Type:** `bool`  
**Default:** `false`

Determines if a bind can be overridden if it's registered twice. When `false`, attempting to register the same type
twice (using either `Register` or `Override`) will panic. When `true`, you can override existing registrations.

> **ℹ️ INFO:** `CanOverride` must be `true` to use the `Override` function. Without it, calling `Override` will panic
> even if you explicitly want to override a dependency.

```go
injector := remy.NewInjector(remy.Config{
    CanOverride: true,
})

// First registration
remy.Register(injector, remy.Instance("first"))

// Using Override function (requires CanOverride: true)
remy.Override(injector, remy.Instance("second"))

// Or using Register again (also works when CanOverride: true)
remy.Register(injector, remy.Instance("third"))

value := remy.MustGet[string](injector)
// value is now "third"
```

**When `CanOverride: false`:**

```go
injector := remy.NewInjector(remy.Config{
    CanOverride: false, // Default
})

remy.Register(injector, remy.Instance("first"))

// This will panic - CanOverride is false
remy.Override(injector, remy.Instance("second")) // ❌ Panics!

// This will also panic - CanOverride is false
remy.Register(injector, remy.Instance("third")) // ❌ Panics!
```

**Use cases:**

- Testing scenarios where you want to replace dependencies
- Development environments where you need flexibility
- Plugin systems where modules can override base dependencies

### DuckTypeElements

**Type:** `bool`  
**Default:** `false`

Enables duck typing for element retrieval. When enabled, Remy can discover and return elements that implement the
requested interface, even if they weren't registered with that exact type.

**How it works:**

When `DuckTypeElements` is enabled, you can use `Get[interfaceName]` directly. Remy will internally call `GetAll` to
search through all registered elements and find matches. If exactly one match is found, it returns that element. If
multiple matches are found, it raises an error.

> **⚠️ CAUTION:** This option is computationally expensive, as it must check all registered elements to find matches.

```go
injector := remy.NewInjector(remy.Config{
    DuckTypeElements: true,
})

// Register a concrete type
remy.Register(injector, remy.Singleton(
    func (retriever remy.DependencyRetriever) (*MyService, error) {
        return &MyService{}, nil
    },
))

// Get can now find MyService even if we request the interface
// (internally uses GetAll to search for matches)
service, err := remy.Get[ServiceInterface](injector)
// service will be *MyService if it implements ServiceInterface

// Or use GetAll to get all matching implementations
services, err := remy.GetAll[ServiceInterface](injector)
// services will contain all services that implement ServiceInterface
```

> **⚠️ WARNING:** If `Get` finds multiple matches, it will return an error:

```go
// Register multiple implementations
remy.Register(injector, remy.Instance(&ServiceA{}))
remy.Register(injector, remy.Instance(&ServiceB{}))

// This will return an error - multiple matches found
service, err := remy.Get[ServiceInterface](injector)
// err: ErrMultipleDIDuckTypingCandidates (2 candidates found)

// Use GetAll instead to get all matches
services, err := remy.GetAll[ServiceInterface](injector)
// services will contain both ServiceA and ServiceB
```

**Use cases:**

- Plugin architectures where multiple implementations exist
- Service discovery patterns
- Testing scenarios with mock implementations
- When you want to retrieve by interface without knowing the concrete type

### ParentInjector

**Type:** `Injector`  
**Default:** `nil`

Creates a child injector that can access all elements registered in the parent injector. The child injector can have its
own additional registrations, but the parent cannot access the child's registrations (scope-safe).

```go
// Create parent injector
parent := remy.NewInjector()
remy.RegisterInstance(parent, "parent-value")

// Create child injector
child := remy.NewInjector(remy.Config{
    ParentInjector: parent,
})

// Child can access parent's registrations
value := remy.MustGet[string](child) // Gets "parent-value"

// Register something in child
remy.RegisterInstance(child, 42)

// Parent cannot access child's registrations
// This would fail:
// remy.Get[int](parent) // Error: element not found
```

**Use cases:**

- Modular applications with scoped dependencies
- Request-scoped injectors in web applications
- Testing with isolated dependency scopes
- Plugin systems with base and extension dependencies

### MultiBinding

**Type:** `bool`  
**Default:** `false`

Enables the `GetAll` family on its own, without turning on duck typing. `GetAll[T]` lists every element registered under
the requested tag that satisfies `T`, and `Get[T]` keeps resolving by exact key.

Before this option existed, the only way to list several implementations was `DuckTypeElements`, which also made every
single `Get` scan the whole registry. `MultiBinding` gives you the list without paying that cost on the hot path.

```go
injector := remy.NewInjector(remy.Config{
    MultiBinding: true,
})

// Each implementation registers under its own concrete type, sharing one tag
remy.RegisterInstance(injector, Base64Encoder{}, "encoders")
remy.RegisterInstance(injector, HexEncoder{}, "encoders")

encoders := remy.MustGetAll[Encoder](injector, "encoders")
// len(encoders) is 2, and no Get paid for a registry scan
```

> **INFO:** Register each implementation under its own concrete type. Two binds sharing a key and a tag collide, the
> same as anywhere else. `GetAll[T]` is what filters the tag down to the ones satisfying `T`.

**Use cases:**

- Plugin lists and handler chains, where every implementation is wanted
- Value groups, when you register several binds under one tag on purpose

### Isolated

**Type:** `bool`  
**Default:** `false`

Stops a lookup miss from falling back to `ParentInjector`. An isolated scope resolves and lists only the binds it owns.

```go
parent := remy.NewInjector()
remy.RegisterInstance(parent, "parent-value")

isolated := remy.NewInjector(remy.Config{
    ParentInjector: parent,
    Isolated:       true,
})

// The parent is still reachable for anything the scope registers itself,
// but a miss no longer walks up to it.
_, err := remy.Get[string](isolated) // error: not registered
```

**Use cases:**

- Test scopes that must not silently inherit a production bind
- Plugin sandboxes that may only see what was handed to them

### ScopeName

**Type:** `string`  
**Default:** `""` (anonymous)

Labels the scope for diagnostics. It shows up when an error reports which scope failed to resolve, which matters once
sub-injectors nest. Anonymous scopes stay the default and cost nothing.

```go
requestScope := remy.NewInjector(remy.Config{
    ParentInjector: appInjector,
    ScopeName:      "http-request",
})
```

### TraceResolution

**Type:** `bool`  
**Default:** `false`

Makes a failed resolution carry the dependency path that led to it, instead of naming only the type that was missing.

It is opt-in because it costs one allocation per failed `Get`. A successful `Get` pays nothing.

```go
injector := remy.NewInjector(remy.Config{
    TraceResolution: true,
})

// Handler needs Service, which needs Repository, which was never registered.
_, err := remy.Get[Handler](injector)

// The message already reads outermost first: Handler -> Service -> Repository
log.Println(err)

var traced *remy.DependencyResolutionError
if errors.As(err, &traced) {
    // Path walks innermost first, so Repository is the first entry.
    // Key carries type identity, not a name, so compare it with NewBindKey.
    for _, entry := range traced.Path() {
        if entry.Key.ID() == remy.NewBindKey[Repository]().ID() {
            log.Println("Repository is the step that failed, tag:", entry.Tag)
        }
    }
}

// The root cause survives the wrapping either way
if errors.Is(err, remy.ErrElementNotRegistered) {
    // ...
}
```

**Use cases:**

- Wiring a large graph, where "element not registered" alone does not say who asked for it
- Test failures that must point at the bind actually missing

## Global Injector

Remy provides a global injector that can be used without explicitly passing an injector instance. Pass `nil` as the
injector parameter to use the global injector.

```go
// Use global injector
remy.Register(nil, remy.Instance("global-value"))
value := remy.MustGet[string](nil)

// Or set a custom global injector
customGlobal := remy.NewInjector(remy.Config{
    CanOverride: true,
})
remy.SetGlobalInjector(customGlobal)
```

**When to use:**

- Simple applications with a single dependency container
- Quick prototyping
- Applications where a single injector is sufficient

## Sub-Injectors

You can create sub-injectors from an existing injector using the `SubInjector` method. This creates a child injector
with optional override configuration.

```go
parent := remy.NewInjector(remy.Config{
    CanOverride: false,
})

// Create sub-injector with override enabled
child := parent.SubInjector(true)

// Child can override parent's binds
remy.Register(child, remy.Instance("child-value"))
```

## Best Practices

1. **Use default configuration** for most cases - Remy's defaults are optimized for common scenarios
2. **Enable `CanOverride`** only when needed (testing, development)
3. **Use `DuckTypeElements` sparingly** - it has performance implications
4. **Leverage `ParentInjector`** for modular applications with clear dependency scopes
5. **Use `NewBindEntry` or `NewBindEntryTagged` in `GetWithPairs`** - The type key is automatically generated from the
   value's type
6. **Create sub-injectors** for request-scoped or test-scoped dependencies

## Configuration Comparison

| Option             | Default | Performance Impact       | Use Case                          |
|--------------------|---------|--------------------------|-----------------------------------|
| `CanOverride`      | `false` | None                     | Testing, development              |
| `DuckTypeElements` | `false` | High                     | Plugin systems, service discovery |
| `MultiBinding`     | `false` | None on `Get`            | Plugin lists, value groups        |
| `ParentInjector`   | `nil`   | Low                      | Scoped dependencies               |
| `Isolated`         | `false` | None                     | Test scopes, plugin sandboxes     |
| `ScopeName`        | `""`    | None                     | Diagnostics on nested scopes      |
| `TraceResolution`  | `false` | One alloc per failed Get | Wiring a large graph              |

> **ℹ️ INFO:** Remy uses zero-width generic types for bindings, providing compile-time type safety without requiring
> reflection. The type key is automatically generated from the value's type when using `NewBindEntry` or
`NewBindEntryTagged`.
