---
title: "Factory"
weight: 3
menu:
  main:
    parent: type-register
    weight: 3
    params:
      icon: "🔄"
---

Factory creates a new instance every time the service is requested. This is useful for services that should not be
shared or need fresh state on each use.

**Key Points:**

- 🔄 New instance created on each request
- 🚫 No caching or reuse
- 🧵 Safe for concurrent use (each call gets its own instance)

## Example

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

type RequestID string

func init() {
	remy.Register(
		nil,
		remy.Factory(
			func(retriever remy.DependencyRetriever) (RequestID, error) {
				// Each call will generate a new unique request ID
				return RequestID(generateUniqueID()), nil
			},
		),
	)
}
```
