package keyopts

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type GenOption uint8

const (
	KeyOptNone              GenOption = 0b0000
	KeyOptGenerifyInterface GenOption = 0b0001
	KeyOptUseReflectionType GenOption = 0b0010
	KeyOptIgnorePointer     GenOption = 0b0100
)

func FromReflectOpts(options types.ReflectionOptions) (result GenOption) {
	if options.UseReflectionType {
		result |= KeyOptUseReflectionType
	}
	if options.GenerifyInterface {
		result |= KeyOptGenerifyInterface
	}
	return
}
