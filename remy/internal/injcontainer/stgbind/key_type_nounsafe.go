//go:build nounsafe

package stgbind

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func (s baseStorage[T]) keyID(key T) T {
	return key
}

type bindKeyID = types.BindKey
