---
title: "Instance"
weight: 4
menu:
  main:
    parent: type-register
    weight: 4
    params:
      icon: "📦"
---

Instance registers an existing value directly without any factory function. This is the simplest form of registration
for values that are already created.

**Key Points:**

- 📝 Registers an existing value directly
- ⚠️ No concurrency protection (read-only recommended)
- 🎯 Best for immutable values or configuration

## Example

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func init() {
	// Register a configuration value
	remy.Register(
		nil,
		remy.Instance("production"),
	)

	// Or use the convenience function
	remy.RegisterInstance(nil, 42)
	remy.RegisterInstance(nil, true)
}
```

> **⚠️ WARNING:** Instance binds have no protection over concurrency, so they're not recommended for structs that perform
operations that modify their attributes. Use Singleton or LazySingleton for concurrent mutable binds.
