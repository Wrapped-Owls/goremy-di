---
title: "Tagged Bindings"
weight: 5
menu:
  main:
    parent: dependency-lifecycle
    params:
      icon: "🏷️"
---

Tags allow you to register and retrieve multiple instances of the same type by providing a unique identifier. This is
useful when you have multiple bindings of the same type and need to distinguish between them.

**Key Points:**

- 🏷️ Multiple bindings of the same type with different tags
- 📝 Tags are optional, and carry the `remy.Tag` type
- 🔍 Use tags when retrieving to get the specific instance you need

---

## The Tag type

A tag is a `remy.Tag`, a defined string type, not a bare `string`. Untyped string constants satisfy it, so a call
written as a literal needs nothing extra:

```go
remy.RegisterInstance(injector, "postgres://billing", "dsn")
value := remy.MustGet[string](injector, "dsn")
```

A `string` **variable** does need the conversion, which is the compiler stopping an arbitrary string from drifting into
a position that means something:

```go
fromConfig := loadDSNTag()
remy.RegisterInstance(injector, dsn, remy.Tag(fromConfig))
```

## Scoping a tag with NewTag

Plain tags share one flat namespace: two packages that both pick `"dsn"` overwrite each other. `NewTag[Scope]` anchors
the name to a type, the way `context.WithValue` keys work, so a package that cannot name the anchor cannot forge the
tag:

```go
type billing struct{}
type shipping struct{}

var (
	billingDSN  = remy.NewTag[billing]("dsn")
	shippingDSN = remy.NewTag[shipping]("dsn")
)

remy.RegisterInstance(injector, "postgres://billing", billingDSN)
remy.RegisterInstance(injector, "postgres://shipping", shippingDSN)

// Same name, different tags, no collision
billingDSN != shippingDSN // true
```

> **CAUTION:** The suffix `NewTag` appends is an opaque per-process id. Never hardcode it, and never persist it or send
> it across processes.

## Registering with Tags

You can register elements with tags using any registration function by providing the tag as the last optional parameter:

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func init() {
	// Register with tag using Register
	remy.Register(
		injector, remy.Instance("production"), "environment",
	)

	// Register another string with a different tag
	remy.Register(
		injector,
		remy.Instance("kafka"),
		"flavour",
	)

	// Register with tag using convenience functions
	remy.RegisterInstance(injector, "development", "env")
	remy.RegisterSingleton(injector, func(retriever remy.DependencyRetriever) (string, error) { return "staging", nil }, "stage")
}

```

---

## Retrieving with Tags

When retrieving, provide the tag as the last optional parameter to get the specific instance:

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	// Retrieve with tag
	env := remy.MustGet[string](injector, "environment")
	// env is "production"

	// Retrieve the other one
	flavour := remy.MustGet[string](injector, "flavour")
	// flavour is "kafka"

	// Using Get with error handling
	stage, err := remy.Get[string](injector, "stage")
	if err != nil {
		// Handle error
	}
}
```

---

## Use Cases

### Configuration Values

```go
remy.RegisterInstance(injector, "production", "environment")
remy.RegisterInstance(injector, "us-east-1", "region")
remy.RegisterInstance(injector, "https://api.example.com", "api-url")

env := remy.MustGet[string](injector, "environment")
region := remy.MustGet[string](injector, "region")
apiURL := remy.MustGet[string](injector, "api-url")
```

### Multiple Implementations

When you have multiple implementations of the same interface:

```go
type Encoder interface {
	Encode([]byte) []byte
}

type Base64Encoder struct{}
func (Base64Encoder) Encode(b []byte) []byte { /* ... */ }

type HexEncoder struct{}
func (HexEncoder) Encode(b []byte) []byte { /* ... */ }

func init() {
	remy.RegisterInstance(injector, Base64Encoder{}, "b64")
	remy.RegisterInstance(injector, HexEncoder{}, "hex")
}

// Later
b64Enc := remy.MustGet[Encoder](injector, "b64")
hexEnc := remy.MustGet[Encoder](injector, "hex")
```

---

## Important Notes

- Tags are **case-sensitive** - `"environment"` and `"Environment"` are different tags
- Tags are **optional** - if you don't provide a tag, the binding is registered without one
- Tags carry the **`remy.Tag`** type: literals pass as they are, a `string` variable needs `remy.Tag(value)`
- Use **`NewTag[Scope]`** when a name could collide with another package's tag
- Retrieving without a tag will get the untagged binding (if it exists)
- You can have both tagged and untagged bindings of the same type
- Tags work with all bind types: `Instance`, `Factory`, `Singleton`, `LazySingleton`
