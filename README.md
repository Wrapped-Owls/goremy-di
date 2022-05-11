# Talaria

A package that helps in creation of golang dependency injections

## About

The Talaria of Mercury (Latin: tālāria) are winged sandals, a symbol of the Greek messenger god Hermes (Roman equivalent
Mercury).

### Inspiration

This package is heavily inpired
by [Modular Dependency Injection](https://modular.flutterando.com.br/docs/flutter_modular/dependency-injection) module

## Instance registration

The strategy for building an instance with its dependencies comprise register all objects in a module and
manufactures them on demand or in single-instance form(singleton). This 'registration' is called **Bind**.

There are a few ways to build a Bind to register object instances:

- _Bind.singleton_: Build an instance only once when the module starts.
- _Bind.lazySingleton_: Build an instance only once when prompted.
- _Bind.factory_: Build an instance on demand.
- _Bind.instance_: Adds an existing instance.
