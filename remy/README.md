# Remy

To install in your Go project, you must use the 1.18 version of the language, and have the GOMODULE enable.

```shell
go get github.com/wrapped-owls/goremy-di
```

## How it works

All the instances and/or closures must be saved somewhere in-memory, so it will be available to be requested later in
the program execution. When the app starts, it must register all dependencies in a given injector, which can be the
default global injector or a custom one. This injector will hold and delegates how to instantiate an object/interface
requested.

### Global Injector

The easiest way to `register` and `retrieve` a bind, is using the **globalInjector**, but this approach comes with a
drawback: To be safer in a multithreading system, during the retrieval access, it uses a RWMutex`, which make the app
execution slower than directing accessing the injector.

To use the global injector, you must set a custom one using the method `SetGlobalInjector`, or let it be created
automatically. After choosing how the global injector will be created, you must pass a _nil_ as the `Injector` parameter
in `Get[T]`/`Register[T]` functions.

### Register bind elements

Use the function `Register` or one of its variants. When using the default `Register` function, a `Bind` generator
function, with a `Binder` function parameter, should be passed as the second parameter. The `Binder` function receive
a `DependencyRetriever` as parameter, which can be used to get another values that was also injected.

```go
package main

// Create an instance of the database connection
func init() {
	remy.Register(
		core.Injector,
		remy.Singleton(func(retriever remy.DependencyRetriever) *sql.DB {
			db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			if err != nil {
				panic(err)
			}
			return db
		}),
	)
}
```

In this example, we register a database connection using the `Register` function, a `Bind` function and the `Binder`
closure. To make it cleaner to read, we can rewrite this using the `RegisterSingleton` or `RegisterInstance` functions:

```go
package main

// Create an instance of the database connection
func init() {
	db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}
	remy.RegisterSingleton(
		core.Injector,
		db,
	)
}
```

#### Bind types

Singleton is not the only bind that can be used, in total is possible to register 4 different types of bind:

- _Singleton_: Build an instance only once when the bind is registered.
- _LazySingleton_: Build an instance only once when it is requested.
- _Factory_: Build an instance on demand.
- _Instance_: Adds an existing instance.

The main difference about **Instance** and **Singleton** is that the singleton bind has a mutex during its access, so is
safely to be used in a multithreading program, while the instance has no protection over it.

#### Using the **DependencyRetriever**

While registering a bind with a `Binder` closure, is possible to retrieve other registered binds by using
the `DependencyRetriever` parameter. So, it can be used in the same way the injector is used to retrieve registered
elements, it is, using the `Get[T]` function.

### Retrieve injected elements

Using the main feature added in Go1.18, we can retrieve all closures/instances using directly the type, instead of a
key. In order to retrieve the element, you must use the function `Get[T]`, which will search the bind in the given
injector, and then return the element in the type T.

As for an example, if we want to retrieve a database connection to create a factory of repositories, we can do this by
calling the Get function:

```go
package main

func init() {
	log.Println("Registering repositories")
	remy.Register(
		core.Injector,
		remy.Factory(func(retriever remy.DependencyRetriever) core.GenericRepository {
			return repositories.NewGenericDbRepository(remy.Get[*sql.DB](retriever))
		}),
	)
}
```

You can also use the registered element in any place of your application, using either the _global_ injector or a _
local_ one.

```go
package main

func main() {
	// Executing create table query
	dbConn := remy.Get[*sql.DB](core.Injector)
	if _, err := dbConn.Exec("CREATE TABLE programming_languages(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}
}
```
