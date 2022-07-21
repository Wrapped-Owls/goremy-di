# Roadmap

- [x] Create the injector using go1.18 generics
- [x] Allow creating sub-injectors
- [x] Create an example directory that shows how the injector can be used
- [ ] Create an injector bind generator package using `//go:generate` comments
- [ ] Refactor the Storage/Injector retrieval to return an `error` instead of a `bool`
- [ ] Remove use of `reflect.TypeOf` in the injector
    - [x] Implement a method to get the type of `element` without using `reflect.TypeOf`
    - [ ] Only use **reflection** if should generify an interface
