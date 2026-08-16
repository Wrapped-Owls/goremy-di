package types

var _ BindEntry = (*InstancePair[struct{}])(nil)

type (
	// BindEntry is an interface used to pass temporary dependencies to GetWithPairs.
	// It encapsulates a value along with its type key and optional tag for registration
	// in a sub-injector during dependency retrieval.
	BindEntry interface {
		// Entry returns the value and its associated BindKey that will be registered
		// in the sub-injector during dependency retrieval.
		Entry() (any, BindKey)
		// Tag returns the optional tag string used for disambiguation when multiple
		// instances of the same type are registered. Returns an empty string if no tag is set.
		Tag() string
	}
	InstancePair[T any] struct {
		// Key is the direct generated BindKey provided by remy default functions
		Key BindKey

		// Value that will be injected directly. (required)
		Value T

		// tag must be used when registering the Value as a bind
		tag string
	}
)

// NewBindPair creates a new InstancePair with the given value and tag.
// The type key is automatically generated from the generic type parameter T.
func NewBindPair[T any](value T, tag string) InstancePair[T] {
	return InstancePair[T]{Value: value, tag: tag, Key: KeyElem[T]{}}
}

func (ip InstancePair[T]) Tag() string {
	return ip.tag
}

func (ip InstancePair[T]) Entry() (any, BindKey) {
	return ip.Value, ip.Key
}

type BindOptions struct {
	Tag          string
	SoftOverride bool
}
