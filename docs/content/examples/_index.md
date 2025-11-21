---
title: "Examples"
weight: 97
---

Complete working examples are available in
the [examples directory](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples) of the GitHub repository.

## Available Examples

- **[basic](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples/basic)** - Basic usage example showing
  singleton registration and retrieval
- **[bindlogger](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples/bindlogger)** - How to inject loggers and
  other utilities
- **[dynamiconstructor](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples/dynamiconstructor)** - Dynamic
  constructor registration examples
- **[guessing_types](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples/guessing_types)** - Type guessing and
  interface examples

Each example is a complete, runnable Go program that demonstrates specific features of Remy DI.

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
	remy.RegisterSingleton(Injector, func() (string, error) {
		return "Hello from Remy!", nil
	})
}

func main() {
	message := remy.MustGet[string](Injector)
	log.Println(message)
}
```

For more detailed examples, visit
the [examples directory on GitHub](https://github.com/Wrapped-Owls/goremy-di/tree/main/examples).
