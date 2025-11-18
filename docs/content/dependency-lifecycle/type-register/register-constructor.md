---
title: "RegisterConstructor"
weight: 5
menu:
  main:
    parent: type-register
    weight: 5
    params:
      icon: "🏗️"
---

The `RegisterConstructor` functions provide a convenient way to register constructor functions without manually writing
the `Binder` function. These functions automatically handle dependency injection for constructor arguments.

## Overview

Instead of writing:

```go
remy.Register(
	injector,
	remy.Singleton(
		func(retriever remy.DependencyRetriever) (*MyService, error) {
			db, err := remy.Get[*sql.DB](retriever)
			if err != nil {
				return nil, err
			}
			return NewMyService(db), nil
		},
	),
)
```

You can simply write:

```go
remy.RegisterConstructorArgs1(
	injector,
	remy.Singleton[*MyService],
	NewMyService,
)
```

## Available Functions

Remy provides constructor registration functions for constructors with 0 to 4 arguments:

### No Arguments

```go
// Without error
remy.RegisterConstructor[T any](
	i Injector,
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() T,
	optTag ...string,
)

// With error
remy.RegisterConstructorErr[T any](
	i Injector,
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() (T, error),
	optTag ...string,
)
```

### One Argument

```go
// Without error
remy.RegisterConstructorArgs1[T, A any](
	i Injector,
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A) T,
	optTag ...string,
)

// With error
remy.RegisterConstructorArgs1Err[T, A any](
	i Injector,
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A) (T, error),
	optTag ...string,
)
```

### Two Arguments

```go
remy.RegisterConstructorArgs2[T, A, B any](...)
remy.RegisterConstructorArgs2Err[T, A, B any](...)
```

### Three Arguments

```go
remy.RegisterConstructorArgs3[T, A, B, C any](...)
remy.RegisterConstructorArgs3Err[T, A, B, C any](...)
```

### Four Arguments

```go
remy.RegisterConstructorArgs4[T, A, B, C, D any](...)
remy.RegisterConstructorArgs4Err[T, A, B, C, D any](...)
```

## Examples

### Example: No Arguments

```go
type Config struct {
	Environment string
}

func NewConfig() *Config {
	return &Config{Environment: "production"}
}

func init() {
	remy.RegisterConstructor(
		nil,
		remy.Singleton[*Config],
		NewConfig,
	)
}
```

### Example: One Argument

```go
type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func init() {
	// Register the database first
	remy.RegisterSingleton(
		nil,
		func(retriever remy.DependencyRetriever) (*sql.DB, error) {
			return sql.Open("sqlite3", ":memory:")
		},
	)

	// Register the service with automatic dependency injection
	remy.RegisterConstructorArgs1(
		nil,
		remy.Singleton[*Service],
		NewService,
	)
}
```

### Example: Two Arguments

```go
type Repository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewRepository(db *sql.DB, logger *log.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func init() {
	remy.RegisterSingleton(nil, func(retriever remy.DependencyRetriever) (*sql.DB, error) {
		return sql.Open("sqlite3", ":memory:")
	})

	remy.RegisterInstance(nil, log.New(os.Stdout, "", 0))

	remy.RegisterConstructorArgs2(
		nil,
		remy.Singleton[*Repository],
		NewRepository,
	)
}
```

### Example: With Error Handling

```go
type Service struct {
	config *Config
}

func NewServiceWithError(config *Config) (*Service, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}
	return &Service{config: config}, nil
}

func init() {
	remy.RegisterInstance(nil, &Config{Environment: "production"})

	remy.RegisterConstructorArgs1Err(
		nil,
		remy.Singleton[*Service],
		NewServiceWithError,
	)
}
```

## Using with Different Bind Types

You can use `RegisterConstructor` functions with any bind type:

```go
// LazySingleton
remy.RegisterConstructorArgs1(
	nil,
	remy.LazySingleton[*Service],
	NewService,
)

// Factory
remy.RegisterConstructorArgs1(
	nil,
	remy.Factory[*Service],
	NewService,
)

// Singleton
remy.RegisterConstructorArgs1(
	nil,
	remy.Singleton[*Service],
	NewService,
)
```

## How It Works

The constructor functions automatically:

1. Retrieve dependencies from the injector based on argument types
2. Call your constructor function with the resolved dependencies
3. Handle errors if dependency retrieval fails
4. Register the result with the specified bind type

This eliminates boilerplate code and makes dependency injection more declarative and easier to read.
