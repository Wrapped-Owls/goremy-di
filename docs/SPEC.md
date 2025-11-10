# Remy DI Spec

All functionalities spec group by kind

## Injectors

All default/official injectors that can be used along with remy

- [remy.NewInjector](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#NewInjector)
- [remy.NewCycleDetectorInjector](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#NewCycleDetectorInjector)

## Bind Type

All Bind types that remy use when registering and retrieving the injected elements

- [remy.Factory](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Factory)
- [remy.Instance](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Instance)
- [remy.LazySingleton](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#LazySingleton)
- [remy.Singleton](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Singleton)

## Registration

Every function that should be used with an Injector to make possible to register and bind all elements easily

- [remy.Register](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Register)
- [remy.Override](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Override)
- [remy.RegisterInstance](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#RegisterInstance)
- [remy.RegisterSingleton](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#RegisterSingleton)

## Retrieve Bind

All functions that can be used to retrieve the injected element

### Functions that return error

- [remy.Get](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Get) - Returns element and error
- [remy.GetAll](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#GetAll) - Returns all matching elements and
  error
- [remy.GetWithPairs](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#GetWithPairs) - Returns element with
  temporary pairs and error
- [remy.GetWith](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#GetWith) - Returns element with callback and
  error

### Functions that panic on error (Must*)

- [remy.MustGet](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MustGet) - Returns element, panics on error
- [remy.MustGetAll](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MustGetAll) - Returns all matching
  elements, panics on error
- [remy.MustGetWithPairs](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MustGetWithPairs) - Returns element
  with temporary pairs, panics on error
- [remy.MustGetWith](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MustGetWith) - Returns element with
  callback, panics on error

### Functions that return zero value on error (Maybe*)

- [remy.MaybeGet](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MaybeGet) - Returns element or zero value on
  error
- [remy.MaybeGetAll](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MaybeGetAll) - Returns all matching
  elements or empty slice on error
- [remy.MaybeGetWithPairs](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MaybeGetWithPairs) - Returns
  element with temporary pairs or zero value on error
- [remy.MaybeGetWith](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#MaybeGetWith) - Returns element with
  callback or zero value on error
