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
    - Add test to check type resolution for elements with same type-name and package-name from another module
- Fix `GetElemKey` method not being able to get the type of the interface
