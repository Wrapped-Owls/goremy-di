---
title: "LazySingleton"
weight: 1
menu:
  main:
    parent: type-register
    params:
      icon: "🚀"
---

LazySingleton creates a single instance of the service, but only when it's first requested. This is useful for expensive
objects that may not always be needed.

**Key Points:**

- 🚀 Lazy services are loaded on first invocation
- 🐎 Lazy service invocation is protected against concurrent loading
- 💾 The instance is cached after first creation

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
		remy.LazySingleton(
			func(retriever remy.DependencyRetriever) (*sql.DB, error) {
				// This will only be called when the database is first requested
				return sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			},
		),
	)
}
```

The database connection will only be established when `Get[*sql.DB]` is called for the first time.
