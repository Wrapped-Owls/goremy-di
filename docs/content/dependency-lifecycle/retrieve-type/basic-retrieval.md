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
db := remy.MustGet[*sql.DB](injector)
// Use db... (will panic if not found)
```

## MaybeGet

`MaybeGet` returns the element or the zero value if an error occurs. Use when you want to ignore errors gracefully.

```go
db := remy.MaybeGet[*sql.DB](injector)
// db will be nil if not found, but no error handling needed
if db != nil {
// Use db...
}
```
