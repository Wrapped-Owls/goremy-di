---
title: "Dependency Lifecycle"
weight: 3
menu:
  main:
    identifier: dependency-lifecycle
    weight: 3
    params:
      icon: "🔄"
---

The dependency lifecycle in Remy DI describes how services are registered, instantiated, and retrieved throughout your
application's execution.

## Design Philosophy

Remy is designed to work seamlessly with your existing code without requiring any modifications to your types or
constructors. Your raw types and constructors work as-is --- no special interfaces to implement, no struct tags to add,
and no wrapper types needed.

**All dependency injection configuration happens only at registration time.** This means:

- ✅ Your existing types and constructors remain unchanged
- ✅ No need to modify your business logic for dependency injection
- ✅ All injection setup is centralized in your registration code
- ✅ Easy to add or remove dependency injection without touching your core code

This approach keeps your codebase clean and maintainable, allowing you to adopt dependency injection incrementally
without refactoring your entire application.

## Next Steps

Learn how to configure your injector:

→ **[Injector Configuration](./injector-configuration/)** - Configure injector behavior, create sub-injectors, and
manage dependency scopes

Learn how to register your dependencies:

→ **[Type Register](./type-register/)** - Discover the different bind types (LazySingleton,
Singleton, Factory, Instance) and how to register your services

Learn how to retrieve your dependencies:

→ **[Retrieve Type](./retrieve-type/)** - Explore the various ways to retrieve registered
dependencies with different error handling strategies
