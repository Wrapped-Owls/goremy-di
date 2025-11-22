---
title: "Retrieving Multiple Elements"
weight: 3
menu:
  main:
    parent: retrieve-type
    params:
      icon: "📦"
---

These functions retrieve all registered elements that match the requested type. This requires `DuckTypeElements` config
to be enabled in the injector configuration.

**Key Points:**

- 📦 Returns all matching implementations
- 🦆 Requires `DuckTypeElements: true` in injector config
- 🔍 Useful for plugin systems and service discovery

---

## Prerequisites

Enable duck typing in the injector configuration:

```go
inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
```

---

## GetAll

`GetAll` retrieves all registered elements that match the requested type and returns an error if something goes wrong.

```go
package main

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

type Service interface {
	DoSomething()
}

func main() {
	// Returns all services that implement the Service interface
	services, err := remy.GetAll[Service](injector)
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services {
		service.DoSomething()
	}
}
```

---

## MustGetAll

`MustGetAll` panics if an error occurs or if `DuckTypeElements` is not enabled.

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	services := remy.MustGetAll[Service](injector)
	// Panics if DuckTypeElements is not enabled or if an error occurs
	for _, service := range services {
		service.DoSomething()
	}
}
```

---

## MaybeGetAll

`MaybeGetAll` returns an empty slice if an error occurs, allowing graceful handling of missing dependencies.

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	services := remy.MaybeGetAll[Service](injector)
	// Returns empty slice if error occurs or DuckTypeElements is not enabled
	for _, service := range services {
		service.DoSomething()
	}
}
```

---

## Use Cases

### Plugin Systems

When you have multiple implementations of the same interface:

```go
package main

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

type Plugin interface {
	Name() string
	Execute()
}

func main() {
	// Register multiple plugins
	remy.RegisterInstance(injector, &PluginA{})
	remy.RegisterInstance(injector, &PluginB{})
	remy.RegisterInstance(injector, &PluginC{})

	// Retrieve all plugins
	plugins := remy.MustGetAll[Plugin](injector)
	for _, plugin := range plugins {
		log.Printf("Executing plugin: %s", plugin.Name())
		plugin.Execute()
	}
}

```

### Service Discovery

Discover all services that implement a specific interface:

```go
package main

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

type HealthCheckable interface {
	HealthCheck() error
}

func main() {
	// All registered health-checkable services
	healthServices := remy.MustGetAll[HealthCheckable](injector)
	for _, service := range healthServices {
		if err := service.HealthCheck(); err != nil {
			log.Printf("Health check failed: %v", err)
		}
	}
}

```

---

> **_Important Notes_**
> - **Ambiguity**: If you need a specific implementation, use tags instead of `GetAll`
> - **Empty Results**: `GetAll` may return an empty slice if no matches are found (this is not an error)
