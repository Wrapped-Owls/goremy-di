package types

type (
	InstancePair[T any] struct {
		// Key is the direct generated BindKey without the automation provided by remy default functions
		Key BindKey

		// Value that will be injected directly. (required)
		Value T

		// Tag must be used when registering the value as a bind
		Tag string
	}
)

type BindOptions struct {
	Tag          string
	SoftOverride bool
}
