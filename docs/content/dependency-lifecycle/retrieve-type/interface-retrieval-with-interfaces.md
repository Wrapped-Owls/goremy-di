---
title: "Retrieving Interfaces (Duck Typing)"
weight: 6
menu:
  main:
    parent: retrieve-type
    params:
      icon: "🦆"
---

# Retrieving Interfaces via Duck Typing

Remy can resolve an **interface** from a **concrete implementation** you registered, as long as duck typing is enabled
in the injector configuration.

This allows you to **register concrete types** and **retrieve them through the interfaces** your consumers depend
on---no explicit interface bindings required.

---

# Prerequisites

Enable duck typing in the injector configuration:

```go
inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
```

---

# Basic Example

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

type Warframe interface{ Cast() string } // Suppose your game logic depends on a Warframe interface

// Citrine is a concrete type implementing Warframe
type Citrine struct{}

func (Citrine) Cast() string { return "Slash Dash" }

func main() {
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})

	// Register the concrete type
	remy.RegisterFactory(inj, remy.Singleton(func() (Citrine, error) {
		return Citrine{}, nil
	}))

	// Retrieve via the interface --- Remy "ducks" to Citrine
	wf := remy.MustGet[Warframe](inj)
	_ = wf.Cast() // => "Slash Dash"
}

```

---

# How Duck Typing Works

When `DuckTypeElements` is enabled:

- Remy resolves services by **assignability**, not only by exact type.
- This allows retrieving an interface (`Warframe`) when only a compatible concrete type (`Citrine`) is registered.
- Lookup follows the scope chain upward; the **first assignable service** is returned.
    - This only works per scope injector, if an Injector instance has multiple possibilities to an interface, it will
      raise an error

---

# Best Practices

## Accept Interfaces, Return Structs --- and Register Structs

The Go proverb _"Accept interfaces, return structs"_ also applies to Remy’s dependency registration:

- **Register concrete types (structs).**
- **Retrieve using interfaces.**

This allows Remy to infer assignability and apply duck typing correctly. If you register using an interface as the key,
Remy cannot guess other interfaces it might also satisfy later. Duck typing helps when a concrete type is registered and
you ask for an interface it implements.

### Recommended Pattern

```go
// Register a concrete implementation
remy.RegisterInstance(inj, Citrine{})

// Retrieve via the interface
_ = remy.MustGet[Warframe](inj)
```

### What _Not_ to Do

```go
// Registering under an interface key hides the concrete type
remy.RegisterInstance[Warframe](inj, Citrine{})

// Later retrieval by another compatible interface will not work
// _ = remy.MustGet[SomeOtherInterface](inj) // fails
```

When you register an interface as the key, Remy can no longer infer what other interfaces the concrete type may satisfy.

---

# Using Tags (Recommended When Multiple Implementations Exist)

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

type Encoder interface{ Encode([]byte) []byte }

type Base64Enc struct{}

func (Base64Enc) Encode(b []byte) []byte { return b }

type HexEnc struct{}

func (HexEnc) Encode(b []byte) []byte { return b }

func main() {
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})

	// Register two implementations with tags
	remy.RegisterInstance(inj, Base64Enc{}, "b64")
	remy.RegisterInstance(inj, HexEnc{}, "hex")

	// Retrieve through the interface using tags
	encB64 := remy.MustGet[Encoder](inj, "b64")
	encHex := remy.MustGet[Encoder](inj, "hex")
}

```

Tags behave normally with duck typing and are often the best way to disambiguate implementations.

---

# Notes & Caveats

### Ambiguous Interfaces

Very broad interfaces (e.g., `fmt.Stringer`) may match many services.
Use tags to avoid ambiguity. When you intentionally want to work with all matching implementations, prefer the
GetAll-family of functions instead of single-value getters.

### Performance

Duck typing broadens the search space.
Enable it only when needed.

### Pointer vs Value Methods

Ensure the concrete type truly satisfies the interface based on receiver type (`value` vs `pointer`).

---

# Troubleshooting

### `ErrElementNotRegistered`

Check that:

- `DuckTypeElements` is enabled.
- The concrete type implements the requested interface.

### Multiple Matches

If more than one implementation matches, you have two main options:

- Add unique tags on registration and use those tags when retrieving a single value, or
- Use the GetAll family to retrieve the full set of matching implementations:

```
// Retrieve all matching services
encoders := remy.MustGetAll[Encoder](inj)
for _, enc := range encoders {
    _ = enc.Encode([]byte("payload"))
}
```

### Child Injectors

Resolution walks up the parent chain; the first assignable service is returned.

---

# Related Configuration

- **`Config.DuckTypeElements`**: enables interface-to-implementation duck typing.
- **`Config.ParentInjector`**: enables parent lookup when using hierarchical injectors.
