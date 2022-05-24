//go:build race

package binds

func (b *SingletonBind[T]) ShouldGenerate() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.dependency == nil
}
