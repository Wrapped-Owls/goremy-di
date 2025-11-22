---
title: "Basic Retrieval"
weight: 2
menu:
  main:
    parent: retrieve-type
    params:
      icon: "🔍"
---

Basic retrieval functions are used to retrieve a single element from the injector.

## Get

`Get` returns the element and an error. Use this when you need to handle errors explicitly.

```go
package main

import (
	"database/sql"
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	db, err := remy.Get[*sql.DB](injector)
	if err != nil {
		log.Fatal(err)
	}
	// Use db...
}
```

## MustGet

`MustGet` returns the element and panics if an error occurs. Use when you're certain the element exists.

```go
package main

import (
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	db := remy.MustGet[*sql.DB](injector)
	// Use db... (will panic if not found)
}

```

## MaybeGet

`MaybeGet` returns the element or the zero value if an error occurs. Use when you want to ignore errors gracefully.

```go
package main

import (
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	db := remy.MaybeGet[*sql.DB](injector)
	// db will be nil if not found, but no error handling needed
	if db != nil {
		// Use db...
	}
}
```

## Retrieving Multiple Matches

When you expect multiple registered elements to match the same requested type (for example, several implementations of
an interface while duck typing is enabled), use the `GetAll` family:

- `GetAll[T]` returns `([]T, error)`
- `MustGetAll[T]` panics on error
- `MaybeGetAll[T]` returns an empty slice on error

```go
// Example: retrieve all encoders implementing the same interface
type Encoder interface{ Encode([]byte) []byte }

encoders := remy.MustGetAll[Encoder](injector)
for _, enc := range encoders {
	_ = enc.Encode([]byte("payload"))
}
```

If you only need one specific implementation among many, prefer registering and retrieving with tags to disambiguate.
