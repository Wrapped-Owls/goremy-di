//go:build go1.24 && !nounsafe

package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type ElementsStorage[T mapKey] struct {
	opts          injopts.CacheConfOption
	namedElements map[string]StorageBackend[uint64, any]
	elements      StorageBackend[uint64, any]
}

func NewElementsStorage[T mapKey](
	opts injopts.CacheConfOption,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
		namedElements: make(map[string]StorageBackend[uint64, any]),
		elements:      newMapBackendUnsafe(),
	}
}

func (s *ElementsStorage[T]) keyID(key T) uint64 {
	return key.ID()
}

func (s *ElementsStorage[T]) newNamedBackend() StorageBackend[uint64, any] {
	return newMapBackendUnsafe()
}

func (s *ElementsStorage[T]) GetAll(optKey ...string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		err = remyErrs.ErrConfigNotAllowReturnAll
		return
	}

	var backend StorageBackend[uint64, any]
	if len(optKey) > 0 && optKey[0] != "" {
		var ok bool
		if b, ok2 := s.namedElements[optKey[0]]; ok2 {
			backend = b
			ok = true
		}
		if !ok {
			return nil, remyErrs.ErrElementNotRegistered{Key: optKey[0]}
		}
	} else {
		backend = s.elements
	}

	resultList = backend.GetAll()
	return
}
