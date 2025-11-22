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
- 📝 Tags are optional string identifiers
- 🔍 Use tags when retrieving to get the specific instance you need

---

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
- Retrieving without a tag will get the untagged binding (if it exists)
- You can have both tagged and untagged bindings of the same type
- Tags work with all bind types: `Instance`, `Factory`, `Singleton`, `LazySingleton`
