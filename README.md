# Remy

A package that helps the creation of golang dependency injections.

## About

The main job of a dependency-injection package is to help providing objects, instances, closures to a caller, by avoiding
a depth graph in parameter pass. Knowing this, and using culinary as inspiration, we thought that a DI (dependency
injection) package is like a `Sous chef` in a cuisine, so we take decided to reference the best `Sous chef` we know: "
Remy"!

### Inspiration

This package is heavily inspired by the 
flutter module [Modular Dependency Injection](https://modular.flutterando.com.br/docs/flutter_modular/dependency-injection)

## Instance registration

The strategy for building an instance with its dependencies comprise on registering all objects in a module and
instantiate them on demand or in single-instance form(singleton). This 'registration' is called **Bind**.

There are a few ways to build a Bind to register object instances:

- _Bind.singleton_: Build an instance only once at the bind registration.
- _Bind.lazySingleton_: Build an instance only once when element is retrieved at first time.
- _Bind.factory_: Build an instance on demand.
- _Bind.instance_: Adds an existing instance as it is.

Other information about binds, functions and injectors can be found easily at [spec](./docs/SPEC.md)
