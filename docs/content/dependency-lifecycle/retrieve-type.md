---
title: "Retrieve Type"
weight: 3
menu:
  main:
    parent: dependency-lifecycle
    params:
      icon: "🔍"
---

Remy DI provides multiple ways to retrieve registered dependencies, each with different error handling strategies. All
retrieval functions use Go generics to provide type-safe dependency resolution.

## Overview

There are three main categories of retrieval functions:

- **Functions without prefix** (`Get`, `GetAll`, `GetWithPairs`, `GetWith`) - Return `(T, error)`
- **Functions with `Must*` prefix** (`MustGet`, `MustGetAll`, etc.) - Panic on error
- **Functions with `Maybe*` prefix** (`MaybeGet`, `MaybeGetAll`, etc.) - Return zero value on error

## Basic Retrieval

### Get

`Get` returns the element and an error. Use this when you need to handle errors explicitly.

```go
package main

import (
	"database/sql"
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	db, err := remy.Get[*sql.DB](injector)
	if err != nil {
		log.Fatal(err)
	}
	// Use db...
}

```

### MustGet

`MustGet` returns the element and panics if an error occurs. Use when you're certain the element exists.

```go
db := remy.MustGet[*sql.DB](injector)
// Use db... (will panic if not found)
```

### MaybeGet

`MaybeGet` returns the element or the zero value if an error occurs. Use when you want to ignore errors gracefully.

```go
db := remy.MaybeGet[*sql.DB](injector)
// db will be nil if not found, but no error handling needed
if db != nil {
// Use db...
}
```

## Retrieving Multiple Elements

### GetAll

`GetAll` retrieves all registered elements that match the requested type (requires `DuckTypeElements` config enabled).

```go
// Returns all services that implement the Service interface
services, err := remy.GetAll[Service](injector)
if err != nil {
log.Fatal(err)
}

for _, service := range services {
service.DoSomething()
}
```

### MustGetAll

```go
services := remy.MustGetAll[Service](injector)
// Panics if DuckTypeElements is not enabled
```

### MaybeGetAll

```go
services := remy.MaybeGetAll[Service](injector)
// Returns empty slice if error occurs
```

## Named Bindings (Tags)

You can register and retrieve elements with tags for named bindings:

```go
// Register with tag
remy.Register(
injector,
remy.Instance("production"),
"environment",
)

// Retrieve with tag
env := remy.MustGet[string](injector, "environment")
```

## Passing Temporary Dependencies

### GetWithPairs

`GetWithPairs` allows you to pass temporary dependencies that will only be available during the retrieval. This is
useful for factory binds that need additional context.

```go
result := remy.MustGetWithPairs[string](
injector,
[]remy.InstancePairAny{
{Value: uint8(42)},
{Value: "Go", Tag: "lang"},
{Value: true},
},
)
```

The temporary pairs are only available to the specific bind being retrieved.

### GetWith

`GetWith` uses a callback function to register temporary dependencies:

```go
result := remy.MustGetWith[string](
injector,
func (injector remy.Injector) error {
remy.RegisterInstance(injector, uint8(42))
remy.RegisterInstance(injector, "Go", "lang")
remy.RegisterInstance(injector, true)
return nil
},
)
```

## Error Handling

### Functions that return error

- `Get[T]` - Returns `(T, error)`
- `GetAll[T]` - Returns `([]T, error)`
- `GetWithPairs[T]` - Returns `(T, error)`
- `GetWith[T]` - Returns `(T, error)`

### Functions that panic (Must\*)

- `MustGet[T]` - Returns `T`, panics on error
- `MustGetAll[T]` - Returns `[]T`, panics on error
- `MustGetWithPairs[T]` - Returns `T`, panics on error
- `MustGetWith[T]` - Returns `T`, panics on error

### Functions that return zero value (Maybe\*)

- `MaybeGet[T]` - Returns `T` or zero value
- `MaybeGetAll[T]` - Returns `[]T` or empty slice
- `MaybeGetWithPairs[T]` - Returns `T` or zero value
- `MaybeGetWith[T]` - Returns `T` or zero value

## Complete Example

```go
package main

import (
	"database/sql"
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

var Injector = remy.NewInjector()

func init() {
	// Register database
	remy.RegisterSingleton(
		Injector,
		func(retriever remy.DependencyRetriever) (*sql.DB, error) {
			return sql.Open("sqlite3", ":memory:")
		},
	)

	// Register with tag
	remy.RegisterInstance(Injector, "production", "environment")
}

func main() {
	// Using Get with error handling
	db, err := remy.Get[*sql.DB](Injector)
	if err != nil {
		log.Fatal(err)
	}

	// Using MustGet (panics on error)
	env := remy.MustGet[string](Injector, "environment")

	// Using MaybeGet (zero value on error)
	config := remy.MaybeGet[Config](Injector)
	if config != nil {
		// Use config...
	}

	log.Printf("Database: %v, Environment: %s", db, env)
}
```

## Best Practices

1. **Use `Get`** when you need explicit error handling
2. **Use `MustGet`** in initialization code where failures should stop the application
3. **Use `MaybeGet`** for optional dependencies
4. **Use tags** for named bindings when you have multiple instances of the same type
5. **Use `GetWithPairs` or `GetWith`** for factory binds that need temporary context
