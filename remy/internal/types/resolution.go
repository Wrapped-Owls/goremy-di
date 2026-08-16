package types

// PathEntry identifies one step of a resolution: the key requested and the
// optional tag used to disambiguate it
type PathEntry struct {
	Key BindKey
	Tag string
}

// GraphNode identifies one dependency in a resolution graph; it is the same
// (key, tag) pair a failed resolution records as its path
type GraphNode = PathEntry
