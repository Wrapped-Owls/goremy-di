---
title: "Named Bindings"
weight: 5
menu:
  main:
    parent: retrieve-type
    params:
      icon: "🏷️"
---

You can register and retrieve elements with tags for named bindings. This is useful when you have multiple instances of
the same type and need to distinguish between them.

## Registering with Tags

```go
// Register with tag
remy.Register(
    injector, remy.Instance("production"), "environment",
)

// Register another string with a different tag
remy.Register(
    injector, remy.Instance("kafka"), "flavour",
)

```

## Retrieving with Tags

```go
// Retrieve with tag
env := remy.MustGet[string](injector, "environment")

// Retrieve the other one
flavour := remy.MustGet[string](injector, "flavour")
```

Tags allow you to have multiple bindings of the same type in the same injector, each identified by a unique tag.
