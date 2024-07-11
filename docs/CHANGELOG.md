# Changelog

## 20220524 - remy/v1.0.0

- Core release
- Register and Retrieve using `generics`
- Example in repository

## 20220530 - remy/v1.1.0

- Create README with more detailed instructions
- Fix some typos
- Add method `SetGlobalInjector`
- Allow to define a **ParentInjector** in `Injector` constructor
- Add more documentation in public types/methods
- Create methods `GetGen` and `GetGenFunc` to pass values dynamically

## 20220721 - remy/v1.2.0

- Replace default type resolution - Now it will not use the `reflect` _package by default_.
- Add `UseReflectionType` option in `Config` struct
- Improve tests coverage
    - Add test to check type resolution for elements with same _type-name_ and _package-name_ from another module
- Fix `GetElemKey` method not being able to get the type of the interface
- Fix an error where `interface` and `pointer` of the same type were being registered as the same type
- Create additional `"Do"` methods: **DoGet**, **DoGetGen**, **DoGetGenFunc**
- Refactor the **Storage/Injector** retrieval to return an `error` instead of a `bool`

## 20220724 - remy/v1.2.1

- Fix an error with `ReflectionOptions` not being applied to **StdInjector**
- Rename some internal attributes in storage

## 20220728 - remy/v1.3.0

- Improve test coverage
- Rename some internal attributes in storage
    - Test `SetGlobalInjector`
    - Test most generics utilities
- Return error on _internal register_ function
- Add `godoc` lines to **internal.types**
- Fix hidden error on `GetGen` function
- Fix error with bind register
    - Prevent overriding a same type with different type of bindings
- Cleanup `DependencyRetriever` methods
    - Removed unnecessary duplicate methods
- Internal improvements
    - Change `BindKey` type to prevent misplace errors
    - Remove duplicate use of **storage**

## 20220801 - remy/v1.4.0

- Create `CycleDetectorInjector` to be used in tests
    - Create a new error type
    - Create a new type in internal utilities
- Change use of unexported type to an exported in public pkg
    - remy public functions now use `Bind[T]` instead of `types.Bind[T]`
- Add `WrapRetriever` to **DependencyRetriever** interface
- Add panic recover to `Do` functions
- Remove `sync.RWMutex` from **globalInjector**
- Swap type `Injector` by `DependencyRetriever` in **Get** methods
- Boost performance by using pointer receiver in _Injector/Storage_ methods

## 20220926 - remy/v1.5.0

- Move package **utils** from `internal` to `pkg`
- Swap all string length comparison with empty string check
- Add `BindKey` & `ReflectOptions` to public API
- Add tests for `func` type
- Update README/SPEC markdown
- Add an error to `Get` method
- Create an example to know how to get logs from Injector
    - This example needs to use unexported packages yet, so it will need update after the release
- Update `bind.Instance` to not use **remy.Binder** functions
- Change `RegisterSingleton` function signature

## 20221006 - remy/v1.6.0

- Change `remy.Binder` to return an error alongside with the value
    - Update all tests to match the new Binder function
- Add new option in `InstancePair` to bind interfaces
- Change the internal use of `ReflectionOptions` to use bitwise operators
    - In this way it'll be possible to add more internal options in the future

## 20230407 - remy/v1.7.0

- Add a new config option to guess element type
    - It only works for instance binds
    - Is strongly recommended to not use
    - Add test cases for new element guessing system
- Create an example that uses the new element guessing option
- Add new errors on utils package
- Remove unused `Default[T any]() T` function from utils
- Refactor cacheConfig internally to use bitwise operators on numeric element
    - This will allow to pass less parameters to constructors

## 20230617 - remy/1.8.0

- Replace internal type detection
    - Stop using `fmt.Sprintf` which uses reflection by default
    - Optimize existent function `TypeNameByReflection`
    - Add **0-width** generic type that will be used as key for injections
- Change `BindKey` type to any
- Upgrade minimal go-version to _1.20_
    - This was made to be able to use any as comparable
- Remove `GenerifyInterfaces=true` from default injector Config

## 20240526 - remy/1.8.2

- Update `GetAll` method (and also duck type get) to retrieve objects stored as **Bind[T]**
    - This update add a little more overhead to check if the bind implements the type we're searching for
    - Created a new function on Binds to try to prevent stack-overflow and cycle dependency get
    - If everything matches, this new way of returning duck-type elements will call the Generate method on bind
        - It must be very careful about cycle dependencies now more than ever
- Add tests to cover the new duckType mode with `Bind[T]`
- Add new functions to register constructors with the injector.
    - `RegisterConstructor` - `RegisterConstructorErr`
    - `RegisterConstructorArgs1` - `RegisterConstructorArgs1Err`
    - `RegisterConstructorArgs2` - `RegisterConstructorArgs2Err`
- Add tests to cover new constructor registration.
