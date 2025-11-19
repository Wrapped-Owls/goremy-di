---
title: "Remy DI"
weight: 1
menu:
  main:
    identifier: home
---

**⚙️ A powerful dependency injection toolkit for Go using generics.**

Remy DI is a modern dependency injection library built for Go 1.18+ that leverages generics to provide a type-safe,
reflection-free API. It implements the Dependency Injection design pattern and can replace other DI packages with a
cleaner, more intuitive interface.

> **See also:**
>
> - [Go package documentation](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy)

## Why Remy DI?

- **🎯 Type-Safe**: Uses Go generics instead of reflection for compile-time type safety
- **🚀 Fast**: No reflection overhead in the hot path
- **📦 Lightweight**: Zero external dependencies
- **🔧 Flexible**: Multiple binding types (Singleton, LazySingleton, Factory, Instance)
- **🌳 Scoped**: Support for sub-injectors and dependency scoping
- **🛡️ Safe**: Built-in circular dependency detection

### **Why this name?**

We love the **short name** for such a utility library. This name references "Remy" - the best Sous chef we know! A DI
package is like a Sous chef in a cuisine, helping to provide objects, instances, and closures to callers.
