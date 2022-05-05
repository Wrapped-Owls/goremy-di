package binds

import "sync"

type SingletonBind[T any] struct {
	dependency T
	mutex      sync.Mutex
}

/*func (b SingletonBind[T]) Get() (T, string) {
	b.mutex.Lock()

}*/
