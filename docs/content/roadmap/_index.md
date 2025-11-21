---
title: "Roadmap"
weight: 99
---

## Core Foundation

- [x] Create the injector using Go 1.18 generics (v1.0.0)
- [x] Allow creating sub-injectors (v1.0.0)
- [x] Refactor the Storage/Injector retrieval to return an `error` instead of a `bool` (v1.2.0, v1.3.0)
- [x] Remove use of `reflect.TypeOf` in the injector (v1.2.0, v1.8.0)
    - [x] Implement a method to get the type of `element` without using `reflect.TypeOf`
    - [x] Only use reflection if explicitly requested by the user
- [x] Create an example directory showing how the injector can be used (v1.0.0)
- [x] Add context-aware resolution (`WithContext`)  (v1.10.0)
    - Works like the current GetWithPairs method, but it will explicitly pass a only the context.Context
- [ ] Create an injector bind generator package using `//go:generate` comments

---

## Registration & Binding

- [x] Add RegisterConstructor wrappers to allow registering constructors directly (v1.8.2)
- [x] Add alias & multi-binding support (v1.8.2)
    - [x] Allow registering a provider under multiple tags
    - [x] Support multiple implementations of the same interface
- [ ] Add RegisterConstructorVariadic wrapper, to allow constructors like NewValue(...T)
- [x] Add module pattern support  (v1.10.0)

---

## Error Handling & Developer Experience

- [ ] Improve error messages
    - [x] Swap error type to include more information about its origin (v1.9.0)
    - [ ] Add dependency path traces (e.g., A → B → C failed)
- [ ] Add dependency graph visualization
    - [ ] Probably export dependency graph to DOT/Graphviz
    - [ ] Add API: `GetDependencyGraph(injector remy.Injector)`
        - Due to the need of building the graph internally, it must be behind a build tag

---

## Advanced Features

- [ ] Add lifecycle hooks
    - [ ] `OnRegister`, `OnResolve` callbacks
    - [ ] Support hook chaining
- [ ] Provide a fluent builder API
    - [ ] Example: `CreateInjector().WithLogger(l).WithModules(m1, m2).Build()`
