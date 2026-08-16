package types

// PathEntry identifies one step of a resolution: the key requested and the
// optional tag used to disambiguate it
type PathEntry struct {
	Key BindKey
	Tag string
}
