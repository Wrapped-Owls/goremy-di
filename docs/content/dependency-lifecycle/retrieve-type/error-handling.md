---
title: "Error Handling"
weight: 6
menu:
  main:
    parent: retrieve-type
    params:
      icon: "⚠️"
---

Remy DI provides three different error handling strategies for retrieval functions. Choose the one that best fits your
use case.

## Functions that Return Error

These functions return `(T, error)` and require explicit error handling:

- `Get[T]` - Returns `(T, error)`
- `GetAll[T]` - Returns `([]T, error)`
- `GetWithPairs[T]` - Returns `(T, error)`
- `GetWith[T]` - Returns `(T, error)`

**Use when:** You need to handle errors explicitly and provide custom error handling logic.

```go
db, err := remy.Get[*sql.DB](injector)
if err != nil {
    // Custom error handling
    log.Printf("Failed to retrieve database: %v", err)
    return err
}
```

## Functions that Panic (Must\*)

These functions return `T` and panic if an error occurs:

- `MustGet[T]` - Returns `T`, panics on error
- `MustGetAll[T]` - Returns `[]T`, panics on error
- `MustGetWithPairs[T]` - Returns `T`, panics on error
- `MustGetWith[T]` - Returns `T`, panics on error

**Use when:** You're certain the element exists and failures should stop the application (e.g., during initialization).

```go
// Will panic if database is not registered
db := remy.MustGet[*sql.DB](injector)
```

## Functions that Return Zero Value (Maybe\*)

These functions return `T` or the zero value if an error occurs:

- `MaybeGet[T]` - Returns `T` or zero value
- `MaybeGetAll[T]` - Returns `[]T` or empty slice
- `MaybeGetWithPairs[T]` - Returns `T` or zero value
- `MaybeGetWith[T]` - Returns `T` or zero value

**Use when:** The dependency is optional and you want to gracefully handle missing dependencies.

```go
// Returns nil if not found, no error handling needed
config := remy.MaybeGet[Config](injector)
if config != nil {
    // Use config...
}
```
