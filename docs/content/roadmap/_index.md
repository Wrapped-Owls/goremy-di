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
- [ ] Create an injector bind generator package using `//go:generate` comments

---

## Registration & Binding

- [x] Add RegisterConstructor wrappers to allow registering constructors directly (v1.8.2)
- [x] Add alias and multi-binding support (v1.8.2)
    - [x] Allow registering a provider under multiple tags
    - [x] Support multiple implementations of the same interface
- [ ] Add RegisterConstructorVariadic wrapper to allow constructors like NewValue(...T)
- [x] Add module pattern support  (v1.10.0)
- [x] Add explicit interface aliasing with `RegisterAs` (v1.13.0)
    - [x] Compile-time checked through a caller-supplied identity conversion
    - [x] O(1) lookup, so duck typing is no longer needed for interface resolution
- [x] Decouple the `GetAll` family from duck typing via `Config.MultiBinding` (v1.13.0)
- [x] Add typed tags with `Tag` and collision-free `NewTag[Scope]` (v1.13.0)

---

## Error Handling & Developer Experience

- [x] Improve error messages
    - [x] Swap the error type to include more information about its origin (v1.9.0)
    - [x] Add dependency path traces (e.g., A → B → C failed) (v1.13.0)
        - Opt-in via `Config.TraceResolution`, so the default failure path keeps its cost
    - [x] Render the cycle path in `ErrCycleDependencyDetected` (v1.13.0)
- [x] Add dependency graph visualization (v1.13.0)
    - [x] Add API: `NewGraphInjector(Config) (Injector, Graph)`
        - An opt-in decorator over the existing `RetrieverFor` seam instead of a build tag,
          unified with the cycle detector: one ordered path both detects cycles and records edges
    - [x] `Graph.ResolveAll` forces every registered bind, listing what could not be built
    - [ ] Export dependency graph to DOT/Graphviz
        - Renderers live outside the core, over the data-only `Graph` interface

---

## Advanced Features & Lifecycle (Framework Layer)

- [ ] Add lifecycle hooks (via external wrapper)
    - [ ] `OnRegister`, `OnResolve` callbacks
    - [ ] Support hook chaining
- [ ] Implement Graceful Shutdown management
    - [ ] Health check probes for registered services
    - [ ] Shutdown sequence management
- [ ] Provide a fluent builder API
    - [ ] Example: `CreateInjector().WithLogger(l).WithModules(m1, m2).Build()`
- [x] Formalize Scoped Injectors (v1.13.0)
    - [x] Name a scope with `Config.ScopeName`; read it back with `Injector.ScopeName`/`Parent`
    - [x] `Config.Isolated` keeps a scope from falling back to its parent on lookups and `GetAll`
    - [ ] Ensure proper resource cleanup for scoped dependencies independent of the parent injector
