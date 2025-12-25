---
title: "Passing Temporary Dependencies"
weight: 5
menu:
  main:
    parent: retrieve-type
    params:
      icon: "⚡"
---

Both `GetWithPairs` and `GetWith` allow you to pass temporary dependencies that will only be available during the
retrieval. This is particularly useful for factory binds that need additional context or runtime values that aren't
registered in the main injector.

## GetWithPairs

`GetWithPairs` allows you to pass temporary dependencies as a slice of `BindEntry`. This is convenient when you
have a fixed set of values to inject.

The type key is automatically generated from the value's type, so you don't need to manually specify it. Use
`remy.NewBindEntry[T](value)` for values without tags, or `remy.NewBindEntryTagged[T](value, tag)` for values that
need a tag.

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.BindEntry{
        remy.NewBindEntry(uint8(42)),
        remy.NewBindEntryTagged("Go", "lang"),
        remy.NewBindEntry(true),
    },
)
```

The temporary pairs are only available to the specific bind being retrieved and are automatically cleaned up after the
retrieval completes.

### Using Tags

You can use tags to disambiguate between multiple instances of the same type:

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.BindEntry{
        remy.NewBindEntryTagged("production", "environment"),
        remy.NewBindEntryTagged("us-east-1", "region"),
        remy.NewBindEntry(42),
    },
)
```

### Using Interface Types

You can also pass interface types directly:

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.BindEntry{
        remy.NewBindEntry[Language](goLang),
        remy.NewBindEntry(uint8(42)),
    },
)
```

## GetWith

`GetWith` uses a callback function to register temporary dependencies. This approach provides more flexibility and
allows you to perform additional setup logic within the callback.

### Basic Usage

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

### Advanced Usage with Error Handling

The callback function can return an error, which will be propagated:

```go
result, err := remy.GetWith[string](
    injector,
    func (injector remy.Injector) error {
        // Register multiple dependencies
        if err := remy.Register(injector, remy.Instance[uint8](42)); err != nil {
            return err
        }

        // Register with tag
        if err := remy.Register(injector, remy.Instance("Go"), "lang"); err != nil {
            return err
        }

        // Register interface type
        if err := remy.Register[Language](
            injector, remy.Instance[Language](goLang),
        ); err != nil {
            return err
        }

        return nil
    },
)
if err != nil {
    log.Fatal(err)
}
```

### Using with Factory Binds

`GetWith` is particularly useful when working with factory binds that need runtime context:

```go
type RequestHandler struct {
    RequestID string
    UserID    int
    DB        *sql.DB
}

// Factory bind that requires temporary context
remy.Register(
    injector,
    remy.Factory(
        func (retriever remy.DependencyRetriever) (*RequestHandler, error) {
            requestID := remy.MustGet[string](retriever, "request-id")
            userID := remy.MustGet[int](retriever, "user-id")
            db := remy.MustGet[*sql.DB](retriever)

            return &RequestHandler{
                RequestID: requestID,
                UserID:    userID,
                DB:        db,
            }, nil
        },
    ),
)

// Later, when handling a request
handler := remy.MustGetWith[*RequestHandler](
    injector,
    func (injector remy.Injector) error {
        remy.RegisterInstance(injector, generateRequestID(), "request-id")
        remy.RegisterInstance(injector, currentUserID, "user-id")
        return nil
    },
)
```

### Conditional Registration

You can use the callback to conditionally register dependencies:

```go
result := remy.MustGetWith[string](
    injector,
    func (injector remy.Injector) error {
        environment := "development"
        if isProduction {
            environment = "production"
        }
        remy.RegisterInstance(injector, environment, "env")
        remy.RegisterInstance(injector, getConfigValue())
        return nil
    },
)
```

## When to Use Each

- **Use `GetWithPairs`** when you have a simple, fixed set of values to pass
- **Use `GetWith`** when you need:
    - More complex registration logic
    - Conditional dependencies
    - Error handling during registration
    - Dynamic value generation

## Important Notes

- Temporary dependencies are **only available** during the retrieval of the specific bind
- They do **not** persist in the main injector after retrieval completes
- They **override** any existing bindings of the same type/tag in the parent injector for the duration of the retrieval
- The type key is automatically generated from the value's type, ensuring type safety at compile time
