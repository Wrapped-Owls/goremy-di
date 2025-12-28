//go:build !go1.24

package stgbind

import remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"

type (
	Node[T any] struct {
		Data  T
		Left  int // -1 == nil
		Right int // -1 == nil
	}
	TreeNodeData[K stgKey] struct {
		KeyValuePair[K, int]
		ResolvedKey uint64
	}
	backendUnsafeTree[K stgKey] struct {
		arena     []Node[TreeNodeData[K]] // all allocated nodes
		valueData []any
	}
)

func (tnd TreeNodeData[K]) Compare(other TreeNodeData[K]) int {
	if tnd.ResolvedKey < other.ResolvedKey {
		return -1
	}
	if tnd.ResolvedKey > other.ResolvedKey {
		return 1
	}
	return 0
}

func newBackend[K stgKey](initialCap uint16) StorageBackend[K, any] {
	t := &backendUnsafeTree[K]{
		arena:     make([]Node[TreeNodeData[K]], 0, initialCap>>1),
		valueData: make([]any, 0, initialCap>>1),
	}
	return t
}

func (t *backendUnsafeTree[K]) createNode(key K, value any) Node[TreeNodeData[K]] {
	t.valueData = append(t.valueData, value)
	return Node[TreeNodeData[K]]{
		Data: TreeNodeData[K]{
			KeyValuePair: KeyValuePair[K, int]{Key: key, Value: len(t.valueData) - 1},
			ResolvedKey:  key.ID(),
		},
		Left:  -1,
		Right: -1,
	}
}

func (t *backendUnsafeTree[K]) Set(key K, value any, allowOverride bool) (triedOverride bool) {
	if len(t.arena) == 0 { // tree is empty
		newNode := t.createNode(key, value)
		t.arena = append(t.arena, newNode)
		return false
	}

	keyID := key.ID()
	var emptyVal int
	curIdx := &emptyVal
	for {
		targetNode := &t.arena[*curIdx]
		targetID := targetNode.Data.ResolvedKey
		var idxPtr *int
		switch {
		case targetID == keyID: // already present
			triedOverride = true
			if !allowOverride {
				// Skip override value
				return true
			}
			t.valueData[targetNode.Data.Value] = value
			return
		case keyID < targetID:
			idxPtr = &targetNode.Left
		case keyID > targetID:
			idxPtr = &targetNode.Right
		}

		if *idxPtr == -1 {
			newNode := t.createNode(key, value)
			t.arena = append(t.arena, newNode)
			*idxPtr = len(t.arena) - 1
			t.arena[*curIdx] = *targetNode
			return triedOverride
		}
		curIdx = idxPtr
	}
}

func (t *backendUnsafeTree[T]) Get(key T) (any, error) {
	var curIdx int
	keyID := key.ID()
	for curIdx != -1 && curIdx < len(t.arena) {
		targetNode := t.arena[curIdx]
		targetID := targetNode.Data.ResolvedKey
		switch {
		case targetID == keyID:
			return t.valueData[targetNode.Data.Value], nil
		case keyID < targetID:
			curIdx = targetNode.Left
		case keyID > targetID:
			curIdx = targetNode.Right
		}
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: key}
}

func (t *backendUnsafeTree[T]) Size() int {
	return len(t.valueData)
}

func (t *backendUnsafeTree[T]) GetAll() []any {
	return t.valueData
}
