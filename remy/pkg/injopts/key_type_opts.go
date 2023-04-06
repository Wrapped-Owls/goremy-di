package injopts

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type KeyGenOption uint8

const (
	KeyOptNone              KeyGenOption = 0b0000
	KeyOptGenerifyInterface KeyGenOption = 0b0001
	KeyOptUseReflectionType KeyGenOption = 0b0010
	KeyOptIgnorePointer     KeyGenOption = 0b0100
)

func KeyOptsFromStruct(options types.ReflectionOptions) (result KeyGenOption) {
	if options.UseReflectionType {
		result |= KeyOptUseReflectionType
	}
	if options.GenerifyInterface {
		result |= KeyOptGenerifyInterface
	}
	return
}
