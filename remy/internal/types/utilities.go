package types

var _ BindEntry = (*InstancePair[struct{}])(nil)

type (
	BindEntry interface {
		Entry() (any, BindKey)
		Tag() string
	}
	InstancePair[T any] struct {
		// Key is the direct generated BindKey without the automation provided by remy default functions
		Key BindKey

		// Value that will be injected directly. (required)
		Value T

		// tag must be used when registering the Value as a bind
		tag string
	}
)

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
