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

- [remy.Get](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#Get)
- [remy.DoGet](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#DoGet)
- [remy.GetGen](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#GetGen)
- [remy.DoGetGen](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#DoGetGen)
- [remy.GetGenFunc](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#GetGenFunc)
- [remy.DoGetGenFunc](https://pkg.go.dev/github.com/wrapped-owls/goremy-di/remy#DoGetGenFunc)
