//go:build !nounsafe

package stgbind

func (s baseStorage[T]) keyID(key T) uint64 {
	return key.ID()
}

type bindKeyID = uint64
