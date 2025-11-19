---
title: "Retrieving Multiple Elements"
weight: 3
menu:
  main:
    parent: retrieve-type
    params:
      icon: "📦"
---

These functions retrieve all registered elements that match the requested type. This requires `DuckTypeElements` config
to be enabled.

## GetAll

`GetAll` retrieves all registered elements that match the requested type.

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

## MustGetAll

`MustGetAll` panics if an error occurs or if `DuckTypeElements` is not enabled.

```go
services := remy.MustGetAll[Service](injector)
// Panics if DuckTypeElements is not enabled
```

## MaybeGetAll

`MaybeGetAll` returns an empty slice if an error occurs.

```go
services := remy.MaybeGetAll[Service](injector)
// Returns empty slice if error occurs
```
