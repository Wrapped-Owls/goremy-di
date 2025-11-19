---
title: "Retrieve Type"
weight: 3
menu:
  main:
    parent: dependency-lifecycle
    identifier: retrieve-type
    params:
      icon: "🔍"
---

Remy DI provides multiple ways to retrieve registered dependencies, each with different error handling strategies. All
retrieval functions use Go generics to provide type-safe dependency resolution.

## Retrieval Methods

There are three main categories of retrieval functions:

- **Functions without prefix** (`Get`, `GetAll`, `GetWithPairs`, `GetWith`) - Return `(T, error)`
- **Functions with `Must*` prefix** (`MustGet`, `MustGetAll`, etc.) - Panic on error
- **Functions with `Maybe*` prefix** (`MaybeGet`, `MaybeGetAll`, etc.) - Return zero value on error

Each category provides the same functionality but with different error handling strategies, allowing you to choose the
approach that best fits your use case.

### Choosing the Right Retrieval Method

1. **Use `Get`** when you need explicit error handling
    - Best for production code where you want to handle errors gracefully
    - Allows you to provide meaningful error messages or fallback behavior
2. **Use `MustGet`** in initialization code where failures should stop the application
    - Use during application startup when missing dependencies indicate a configuration error
    - Prevents the application from running in an invalid state
3. **Use `MaybeGet`** for optional dependencies
    - Use when a dependency might not be available and that's acceptable
    - Simplifies code by avoiding error handling for optional features

### Using Tags

Use tags for named bindings when you have multiple instances of the same type

- Helps distinguish between different configurations or implementations
- Makes code more readable and maintainable
