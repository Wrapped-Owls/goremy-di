package remy

import (
	"strconv"

	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

// Tag disambiguates multiple binds of the same type. Untyped string constants
// still satisfy it; string variables need an explicit Tag(value) conversion.
type Tag string

// NewTag creates a Tag namespaced by the Scope anchor type, like context.WithValue
// keys: a package that cannot name the anchor cannot forge the tag.
//
// The suffix is an opaque per-process id, so it must not be hardcoded or persisted.
func NewTag[Scope any](name string) Tag {
	return Tag(name + "@" + strconv.FormatUint(utils.NewKeyElem[Scope]().ID(), 16))
}
