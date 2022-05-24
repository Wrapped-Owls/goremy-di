//go:build !race

package binds

func (b *SingletonBind[T]) ShouldGenerate() bool {
	return b.dependency == nil
}
