---
title: "Getting Started"
weight: 2
---

## Installation

Install Remy DI in your Go project (requires Go 1.20+):

```shell
go get github.com/wrapped-owls/goremy-di/remy
```

## Quick Example

Here's a simple example to get you started:

```go
package main

import (
	"log"
	"github.com/wrapped-owls/goremy-di/remy"
)

var Injector = remy.NewInjector()

func init() {
	// Register a service as singleton
	remy.RegisterSingleton(Injector, func() (string, error) {
		return "Hello, Remy!", nil
	})
}

func main() {
	// Retrieve the service
	message := remy.MustGet[string](Injector)
	log.Println(message) // Output: Hello, Remy!
}
```

## Key Concepts

- **Injector**: The container that manages your dependencies. You can create multiple injectors or use the global one.
- **Registration**: Register your services during application startup using one of the bind types.
- **Retrieval**: Get your registered services anywhere in your code using type-safe generic functions.
