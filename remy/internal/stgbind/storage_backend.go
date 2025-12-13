package stgbind

// StorageBackend defines the interface for storage backends (tree or map)
type StorageBackend[K comparable, V any] interface {
	// Set stores a value with the given key
	// allowOverride indicates if overriding an existing key is allowed
	// Returns triedOverride=true if the key already existed
	Set(key K, value V, allowOverride bool) (triedOverride bool)
	// Get retrieves a value by key
	Get(key K) (V, error)
	// Size returns the number of elements in the storage
	Size() int
	// GetAll returns all values in the storage
	GetAll() []V
}
