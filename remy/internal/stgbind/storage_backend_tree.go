//go:build !go1.24

package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

// uint64KeyValuePair stores a uint64 key and its value for the tree
type uint64KeyValuePair struct {
	Key   uint64
	Value any
}

// treeBackend implements StorageBackend[uint64, any] using a binary tree
type treeBackend struct {
	tree *Tree[uint64KeyValuePair]
}

// newTreeBackend creates a new tree backend
func newTreeBackend() StorageBackend[uint64, any] {
	return &treeBackend{
		tree: NewBinaryTree[uint64KeyValuePair](func(a, b uint64KeyValuePair) int {
			if a.Key < b.Key {
				return -1
			}
			if a.Key > b.Key {
				return 1
			}
			return 0
		}, 11),
	}
}

func (tb *treeBackend) Set(key uint64, value any, allowOverride bool) (triedOverride bool) {
	pair := uint64KeyValuePair{Key: key, Value: value}
	searchPair := uint64KeyValuePair{Key: key}
	idx := tb.tree.Search(searchPair)

	if idx != -1 {
		triedOverride = true
		if !allowOverride {
			return true
		}
		// Update existing node
		tb.tree.Update(idx, pair)
		return true
	}

	// Insert new node
	_, _ = tb.tree.Insert(pair)
	return false
}

func (tb *treeBackend) Get(key uint64) (any, error) {
	searchPair := uint64KeyValuePair{Key: key}
	idx := tb.tree.Search(searchPair)
	if idx == -1 {
		return nil, remyErrs.ErrElementNotRegistered{Key: key}
	}
	pair := tb.tree.GetValue(idx)
	return pair.Value, nil
}

func (tb *treeBackend) Size() int {
	return tb.tree.Size()
}

func (tb *treeBackend) GetAll() []any {
	result := make([]any, 0, tb.tree.Size())
	tb.tree.InOrder(func(pair uint64KeyValuePair) {
		result = append(result, pair.Value)
	})
	return result
}
