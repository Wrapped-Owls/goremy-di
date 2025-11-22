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

`GetWithPairs` allows you to pass temporary dependencies as a slice of `InstancePairAny`. This is convenient when you
have a fixed set of values to inject.

> **Important:** When using `GetWithPairs`, you must specify the bind key using `remy.NewBindKey[T]()` for each value.
> This is required to work without reflection enabled. If you omit the `Key` field, the injector will attempt to use
> reflection to determine the type, which requires the `UseReflectionType` configuration to be enabled.

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.InstancePairAny{
        {Key: remy.NewBindKey[uint8](), Value: uint8(42)},
        {Key: remy.NewBindKey[string](), Value: "Go", Tag: "lang"},
        {Key: remy.NewBindKey[bool](), Value: true},
    },
)
```

The temporary pairs are only available to the specific bind being retrieved and are automatically cleaned up after the
retrieval completes.

### Using Interface Values

You can also specify interface types explicitly using `InterfaceValue`, which will require reflection:

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.InstancePairAny{
        {
            Key:            remy.NewBindKey[Language](),
            Value:          goLang,
            InterfaceValue: (*Language)(nil), // Explicitly bind as Language interface
        },
    },
)
```

### Direct Bind Keys

For more control, you can specify the bind key directly:

```go
result := remy.MustGetWithPairs[string](
    injector,
    []remy.InstancePairAny{
        {Key: remy.NewBindKey[uint8](), Value: uint8(42)},
        {Key: remy.NewBindKey[string](), Value: "Go", Tag: "lang"},
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
            injector,
            remy.Instance[Language](goLang),
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
- The sub-injector created for temporary dependencies inherits the reflection options from the parent injector
