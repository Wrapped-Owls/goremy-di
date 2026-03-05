# Remy DI

Type-safe dependency injection for Go using generics.

Remy DI provides a lightweight and fast way to register and resolve dependencies without reflection in the hot path.
It supports scoped injectors, bind tags, factory parameters, and circular dependency detection.

## About

The main job of a dependency-injection package is to help providing objects, instances, closures to a caller, by
avoiding
a depth graph in parameter pass. Knowing this, and using culinary as inspiration, we thought that a DI (dependency
injection) package is like a `Sous chef` in a cuisine, so we take decided to reference the best `Sous chef` we know: "
Remy"!

### Inspiration

This package is heavily inspired by the
flutter
module [Modular Dependency Injection](https://modular.flutterando.com.br/docs/flutter_modular/dependency-injection)

## Documentation

- Docs site: https://wrapped-owls.github.io/goremy-di/
- Go package: https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy
- Benchmarks: [benchmark/README.md](./benchmark/README.md)

## Requirements

- Go `1.20+`

## Installation

```bash
go get github.com/wrapped-owls/goremy-di/remy
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

var Injector = remy.NewInjector()

func init() {
	remy.RegisterSingleton(Injector, func(_ remy.DependencyRetriever) (string, error) {
		return "hello from remy", nil
	})
}

func main() {
	message := remy.MustGet[string](Injector)
	log.Println(message)
}
```

## Core Concepts

- `Injector`: Container that stores binds and resolves dependencies.
- `Bind`: Registration strategy for a type.
- `DependencyRetriever`: Resolver interface used inside bind constructors.
- Optional tags: Distinguish multiple binds with the same type.
- Sub-injectors: Child scopes can inherit from a parent injector.

### Bind Types

- `Singleton`: Constructed once when registered.
- `LazySingleton`: Constructed once on first retrieval.
- `Factory`: Constructed on every retrieval.
- `Instance`: Existing value registered as-is.

## Project Structure

- `remy/`: Main library module (`go.mod`, implementation, tests).
- `examples/`: Runnable usage samples (`basic`, `bindlogger`, `dynamiconstructor`, `guessing_types`).
- `docs/`: Hugo documentation source and generated site assets.
- `benchmark/`: Benchmark suites and benchstat comparison output.

## Development Commands

Run from repository root:

```bash
make test
make test-race
make run-lint
make run-formatters
make examples
make docs-build
make docs-serve
make benchmark
```

## Notes

- Prefer registering dependencies during startup, before concurrent retrieval.
- Use `Get[T]` if you want explicit error handling, or `MustGet[T]` when failure should panic.
- For advanced usage (modules, tags, temporary dependencies, scoped retrieval), see the docs site.
