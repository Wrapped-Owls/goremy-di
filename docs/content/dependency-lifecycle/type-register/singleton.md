---
title: "Singleton"
weight: 2
menu:
  main:
    parent: type-register
    params:
      icon: "⚡"
---

Singleton creates a single instance immediately when the bind is registered. The instance is built during the
registration phase and reused throughout the application lifetime.

**Key Points:**

- ⚡ Instance is created during registration
- 🔒 Thread-safe singleton pattern
- 📦 Same instance returned for all requests

## Example

```go
package main

import (
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"
)

func init() {
	remy.Register(
		nil,
		remy.Singleton(
			func(retriever remy.DependencyRetriever) (*sql.DB, error) {
				// This is called immediately during registration
				return sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			},
		),
	)
}
```

You can also use the convenience function:

```go
remy.RegisterSingleton(
	nil,
	func(retriever remy.DependencyRetriever) (*sql.DB, error) {
		return sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
	},
)
```
