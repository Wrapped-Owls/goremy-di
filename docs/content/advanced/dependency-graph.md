---
title: "Dependency Graph"
weight: 6
menu:
  main:
    parent: advanced
    params:
      icon: "🕸️"
---

A dependency cycle is hard to see in code and brutal at runtime: the resolution recurses until the stack dies. Remy
gives you two injectors to catch it, both wrappers around a normal one, and both far slower than the standard injector.
Use them in tests and while debugging, never on a hot path.

---

## Detecting cycles

`NewCycleDetector` records the chain it is walking and reports the moment a key repeats inside it. Nothing else is kept.

```go
package main

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestNoCycles(t *testing.T) {
	injector := remy.NewCycleDetector()
	registerEverything(injector)

	if _, err := remy.Get[Handler](injector); err != nil {
		t.Error(err)
	}
}
```

A cycle comes back as an error naming the chain that closed it:

```
cycle dependency detected: types.KeyElem[main.Handler] -> types.KeyElem[main.Service] -> types.KeyElem[main.Handler]
```

The detection follows the resolution into temporary scopes as well, so a cycle that only closes inside a `GetWith` or
`GetWithPairs` binder, or across a `SubInjector` created mid-resolution, is reported instead of hanging.

> **CAUTION:** Every `Get` inside a binder must go through the `DependencyRetriever` the binder received. If a binder
> closes over the injector variable instead, the detector never sees that edge.

---

## Recording the graph

`NewGraphInjector` returns the injector plus a `Graph` view of what it recorded, so you can read the wiring back
instead of guessing it.

```go
injector, dependencies := remy.NewGraphInjector()
registerEverything(injector)

_ = remy.MustGet[Handler](injector)

for _, edge := range dependencies.Edges() {
	// edge.From requested edge.To during a resolution
	log.Println(edge.From, "->", edge.To)
}
```

An edge is only recorded for a resolution that actually happened. To exercise everything that was registered, ask the
graph to build it all:

```go
failed, err := dependencies.ResolveAll()
if err != nil {
	for _, node := range failed {
		log.Printf("%v could not be built: %v", node.Node, node.Err)
	}
}
```

`ResolveAll` returns one `FailedNode` per bind it could not build, each carrying the reason that bind gave, and an
`err` joining all of them. A registration that was never requested by anything still gets built here, which is what
makes it useful as a wiring smoke test.

---

## Which one to use

| Need                                     | Use                  |
|------------------------------------------|----------------------|
| Assert a test wiring has no cycles       | `NewCycleDetector`   |
| Read which bind asked for which          | `NewGraphInjector`   |
| Prove every registration can be built    | `Graph.ResolveAll`   |
| Anything in production                   | `NewInjector`        |

> **INFO:** `NewCycleDetectorInjector` still works and still detects cycles, but it is deprecated. Use
> `NewCycleDetector`, or `NewGraphInjector` when the recorded graph is also wanted.
