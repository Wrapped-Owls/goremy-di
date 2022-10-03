# Remy

To install in your Go project, you must use the 1.18 version of the language, and have the environment
variable `GO111MODULE=on`.

```shell
go get github.com/wrapped-owls/goremy-di/remy
```

## How it works

All the instances and/or closures must be saved somewhere in-memory, so it will be available to be requested later in
the program execution. When the app starts, it must register all dependencies in a given injector, which can be the
default global injector or a custom one. This injector will hold and delegate how to instantiate an object/interface
requested.

### The Injector

Remy can generate multiple injector instances, and by default it can configure a _default_ global instance.
To generate a new Injector instance, the function `NewInjector` must be called passing a `Config` struct, which have the
attributes:

- AllowOverride - Determines if a element bind can be override
- GenerifyInterface - Go interfaces are not exactly unique, so when an element is registered with an interface `A` and
  it has the same methods signatures as interface `B`, this flag will tell to remy to treat both as the same.
- UseReflectionType - Use reflection to get the type of the object/interface. This is useful when your program has a
  type with the same package and type name from another module or subpackage.
- ParentInjector - Make possible to pass an existing Injector as a parent injector, so the new injector can access all
  elements in its parent. Is good to know, that all binds registered in sub-injector can't be accessed by the parent
  injector, it is scope safe.

```go
package core

import (
	"github.com/wrapped-owls/goremy-di/remy"
	"log"
)

var Injector remy.Injector

// create a new instance of the injector
func init() {
	log.Println("Initializing injector")
	Injector = remy.NewInjector()
}
```

### Global Injector

The easiest way to `register` and `retrieve` a bind, is using the **globalInjector**, but this approach comes with a
drawback: To be safer in a multithreading system, during the retrieval access, it uses a RWMutex`, which make the app
execution slower than directing accessing the injector.

A global injector can be defined in two different ways, bt setting a custom one using the method `SetGlobalInjector`, or
let it be created automatically by the **remy** package.

**Curiosity tip:** To don't allocate the object in memory directly, the global injector is only generated when it is
first accessed.

To use the global injector, you must pass a _nil_ as the `Injector` parameter
in `Get[T]`/`Register[T]`/`GetGen[T]` functions.

### Register bind elements

Use the function `Register` or one of its variants. When using the default `Register` function, a `Bind` generator
function, with a `Binder` function parameter, should be passed as the second parameter. The `Binder` function receive
a `DependencyRetriever` as parameter, which can be used to get another values that was also injected.

```go
package main

import (
	"database/sql"
	"github.com/wrapped-owls/goremy-di/remy"
)

// Create an instance of the database connection
func init() {
	remy.Register(
		core.Injector,
		remy.Singleton(
			func(retriever remy.DependencyRetriever) *sql.DB {
				db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
				if err != nil {
					panic(err)
				}
				return db
			},
		),
	)
}
```

In this example, we register a database connection using the `Register` function, a `Bind` function and the `Binder`
closure. To make it cleaner to read, we can rewrite this using the `RegisterSingleton` or `RegisterInstance` functions:

```go
package main

import (
	"database/sql"
	"github.com/wrapped-owls/goremy-di/remy"
)

// Create an instance of the database connection
func init() {
	remy.RegisterSingleton(
		core.Injector, func(retriever remy.DependencyRetriever) *sql.DB {
			db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			if err != nil {
				panic(err)
			}
			return db
		},
	)
}
```

#### Bind types

Singleton is not the only bind that can be used, in total is possible to register 4 different types of bind:

- _Singleton_: Build an instance only once when the bind is registered.
- _LazySingleton_: Build an instance only once when it is requested.
- _Factory_: Build an instance on demand.
- _Instance_: Adds an existing instance.

The main difference about **Instance** and **Singleton** is that the singleton bind was made to be used when the object
needs to be generated using the value in other binds, in this way we still have an instance, but it will be more easily
to generate and inject.

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

import (
	"database/sql"
	"github.com/wrapped-owls/goremy-di/remy"
)

func init() {
	remy.Register(
		core.Injector,
		remy.Factory(
			func(retriever remy.DependencyRetriever) core.GenericRepository {
				return repositories.NewGenericDbRepository(remy.Get[*sql.DB](retriever))
			},
		),
	)
}
```

You can also use the registered element in any place of your application, using either the _global_ injector or a _
local_ one.

```go
package main

import (
	"database/sql"
	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	// Executing create table query
	dbConn := remy.Get[*sql.DB](core.Injector)
	if _, err := dbConn.Exec("CREATE TABLE programming_languages(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}
}
```

#### Passing parameters to bind injector

Sometimes you may want to generate a bind that must receive some additional arguments in the `Binder` function, but
these arguments are variables and depend on the caller. So you may think that the solution for this is to dynamically
during the program runtime, call the `Register[T]` function and register these parameters, but this thought is wrong,
because the bind registration is not _concurrent_ safe. This means that while the register of a single _string_, another
thread can register the same element, and then the bind will retrieve the values in a wrong way.

This works by creating automatically a sub-injector that will be used to add instance binds that will be used to
generate the factory bind requested.

**REMINDER:** It only works with the `Factory` and `LazySingleton` binds.

Currently, exists two ways to do this, by using an array of `InstancePair` or by using a function and registering the
values directly in it.

Using as example a factory bind registered in the init function:

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

func init() {
	remy.Register(
		nil, remy.Factory(
			func(injector remy.DependencyRetriever) string {
				return fmt.Sprintf(
					"I love %s, yes this is %v, as the answer %d",
					remy.Get[string](injector, "lang"), remy.Get[bool](injector), remy.Get[uint8](injector),
				)
			},
		),
	)
}
```

The requested values can be passed by two forms:

##### Using InstancePair array

With this method is not possible to register correctly interfaces, so in case the factory binds requests an interface
value, is better to use the other method.

```go
package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
	"log"
)

func main() {
	result := remy.GetGen[string](
		injector,
		[]remy.InstancePairAny{
			{
				Value: uint8(42),
			},
			{
				Value: "Go",
				Key:   "lang",
			},
			{
				Value: true,
			},
		},
	)

	log.Println(result)
}
```

##### Using callback to register the values

```go
package main

import "github.com/wrapped-owls/goremy-di/remy"

func main() {
	remy.GetGenFunc[string](
		injector, func(injector remy.Injector) {
			remy.Register(ij, remy.Instance[uint8](42))
			remy.Register(ij, remy.Instance("Go"), "lang")
			remy.Register(ij, remy.Instance(true))
		},
	)
}
```

### Cycle dependencies problem

As you can use the binds in a dynamic way, by generating factories that will run during runtime, it may happen to a
bind request some other bind that depends on the first bind, this is a dependency cycle. The main problem with a cycle
like this, is that is really hard to detect during code/compile time, and once the code is running, it can end in
a `Stack Overflow`, which causes a panic to the program, and may be harm to it.

So, to enable the possibility to test and detect for a _dependency-cycle_, you can use the `CycleDetectorInjector`,
which can be called using the constructor **"NewCycleDetectorInjector"**. The main question during it's use is that it
creates a wrap in the `StandardInjector`, and create an internal graph for each dependency that was requested to the
injector. This functionality is much slower than using the `StandardInjector`, so it is only recommended to use it in
test files, to make sure that no dependency cycle was created.

```go

package main

import (
	"github.com/wrapped-owls/goremy-di/remy"
	"testing"
)

func createInjections(injector remy.Injector) {
	// ...
}

func TestCycles(t *testing.T) {
	ij := remy.NewCycleDetectorInjector()
	createInjections(ij)
	if _, err := remy.DoGet[string](ij); err != nil {
		t.Error(err)
	}
}
```

#### Important Note

When using the `CycleDetectorInjector` is important that in Binds, all _Get_ methods used call the
given `DependencyRetriever`, if the same injector is used inside the function, as a clojure, it will not be able to
detect cycles.
