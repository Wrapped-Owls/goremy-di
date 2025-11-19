---
title: "Type Register"
weight: 1
menu:
  main:
    parent: dependency-lifecycle
    identifier: type-register
    params:
      icon: "📝"
---

Remy DI supports four different bind types for registering dependencies. Each type has different instantiation behavior
and lifecycle characteristics.

## Bind Types

- **[LazySingleton](./lazy-singleton/)** - Creates a single instance when first
  requested. Useful for expensive objects that may not always be needed.
- **[Singleton](./singleton/)** - Creates a single instance immediately during
  registration. Thread-safe singleton pattern.
- **[Factory](./factory/)** - Creates a new instance every time the service is
  requested. Safe for concurrent use.
- **[Instance](./instance/)** - Registers an existing value directly. Best for
  immutable values or configuration.

## Advanced Registration

- **[RegisterConstructor](./register-constructor/)** - Convenient way to register
  constructor functions without manually writing the `Binder` function. Automatically handles dependency injection for
  constructor arguments.
