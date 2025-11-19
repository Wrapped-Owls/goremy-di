---
title: "Complete Example"
weight: 7
menu:
  main:
    parent: retrieve-type
    params:
      icon: "💡"
---

Here's a complete example demonstrating various retrieval methods:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

var Injector = remy.NewInjector()

func init() {
	// Register database
	remy.RegisterSingleton(
		Injector,
		func(retriever remy.DependencyRetriever) (*sql.DB, error) {
			return sql.Open("sqlite3", ":memory:")
		},
	)

	// Register with tag
	remy.RegisterInstance(Injector, "production", "environment")

	// Register a factory that needs temporary context
	remy.Register(
		Injector,
		remy.Factory(
			func(retriever remy.DependencyRetriever) (string, error) {
				lang := remy.MustGet[string](retriever, "lang")
				answer := remy.MustGet[uint8](retriever)
				return fmt.Sprintf("I love %s, answer is %d", lang, answer), nil
			},
		),
	)
}

func main() {
	// Using Get with error handling
	db, err := remy.Get[*sql.DB](Injector)
	if err != nil {
		log.Fatal(err)
	}

	// Using MustGet (panics on error)
	env := remy.MustGet[string](Injector, "environment")

	// Using MaybeGet (zero value on error)
	config := remy.MaybeGet[Config](Injector)
	if config != nil {
		// Use config...
	}

	// Using GetWithPairs for temporary dependencies
	result := remy.MustGetWithPairs[string](
		Injector,
		[]remy.InstancePairAny{
			{Key: remy.NewBindKey[uint8](), Value: uint8(42)},
			{Key: remy.NewBindKey[string](), Value: "Go", Tag: "lang"},
		},
	)
	log.Println(result)

	// Using GetWith for more complex temporary dependencies
	result2 := remy.MustGetWith[string](
		Injector,
		func(injector remy.Injector) error {
			remy.RegisterInstance(injector, uint8(100))
			remy.RegisterInstance(injector, "Rust", "lang")
			return nil
		},
	)
	log.Println(result2)

	log.Printf("Database: %v, Environment: %s", db, env)
}
```
