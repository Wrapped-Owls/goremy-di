package stgbind

import "errors"

type Node[T any] struct {
	Value T
	Left  int // -1 == nil
	Right int // -1 == nil
}

// Tree is a binary tree that stores all nodes in a slice.
type Tree[T any] struct {
	arena []Node[T]        // all allocated nodes
	free  []int            // indices that are free
	root  int              // index of the root node (-1 = empty)
	cmp   func(a, b T) int // comparison function: <0 if a<b, 0 if a==b, >0 if a>b
}

// NewBinaryTree creates a tree with a given capacity.  It pre‑allocates the arena
// so that the first `cap` inserts will not trigger a slice reallocation.
func NewBinaryTree[T any](cmp func(a, b T) int, cap int) *Tree[T] {
	t := &Tree[T]{
		arena: make([]Node[T], cap),
		free:  make([]int, cap),
		root:  -1,
		cmp:   cmp,
	}
	// fill free list
	for i := 0; i < cap; i++ {
		t.free[i] = i
	}
	return t
}

// nextIdx pops an index from the free stack or returns -1 if full.
func (t *Tree[T]) nextIdx() int {
	n := len(t.free)
	if n == 0 {
		return -1
	}
	idx := t.free[n-1]
	t.free = t.free[:n-1]
	return idx
}

// ---------- Core operations ----------

// Insert adds value v to the tree.  If v already exists, nothing changes
// and the index of the existing node is returned.
func (t *Tree[T]) Insert(v T) (int, error) {
	if t.root == -1 { // tree empty
		idx := t.nextIdx()
		if idx == -1 {
			return -1, errors.New("tree is full")
		}
		t.arena[idx] = Node[T]{Value: v, Left: -1, Right: -1}
		t.root = idx
		return idx, nil
	}

	cur := t.root
	for {
		cmp := t.cmp(t.arena[cur].Value, v)
		var idxPtr *int
		switch {
		case cmp == 0: // already present
			return cur, nil
		case cmp > 0: // v < current, go left
			idxPtr = &t.arena[cur].Left
		case cmp < 0: // v > current, go right
			idxPtr = &t.arena[cur].Right
		}

		if *idxPtr == -1 {
			idx := t.nextIdx()
			if idx == -1 {
				return -1, errors.New("tree is full")
			}
			t.arena[idx] = Node[T]{Value: v, Left: -1, Right: -1}
			*idxPtr = idx
			return idx, nil
		}
		cur = *idxPtr
	}
}

// Search returns the index of a node that holds v, or -1 if not found.
func (t *Tree[T]) Search(v T) int {
	cur := t.root
	for cur != -1 {
		cmp := t.cmp(t.arena[cur].Value, v)
		if cmp == 0 {
			return cur
		}
		if cmp > 0 {
			cur = t.arena[cur].Left
		} else {
			cur = t.arena[cur].Right
		}
	}
	return -1
}

// GetValue returns the value stored at the given index.
func (t *Tree[T]) GetValue(idx int) T {
	if idx == -1 || idx >= len(t.arena) {
		var zero T
		return zero
	}
	return t.arena[idx].Value
}

// Update updates the value at the given index. Returns false if index is invalid.
func (t *Tree[T]) Update(idx int, v T) bool {
	if idx == -1 || idx >= len(t.arena) {
		return false
	}
	t.arena[idx].Value = v
	return true
}

func (t *Tree[T]) Size() int {
	if t.root == -1 {
		return 0
	}
	var sz int
	var stack []int = []int{t.root}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		sz++
		if t.arena[cur].Right != -1 {
			stack = append(stack, t.arena[cur].Right)
		}
		if t.arena[cur].Left != -1 {
			stack = append(stack, t.arena[cur].Left)
		}
	}
	return sz
}

func (t *Tree[T]) InOrder(f func(v T)) {
	var walk func(idx int)
	walk = func(idx int) {
		if idx == -1 {
			return
		}
		walk(t.arena[idx].Left)
		f(t.arena[idx].Value)
		walk(t.arena[idx].Right)
	}
	walk(t.root)
}
