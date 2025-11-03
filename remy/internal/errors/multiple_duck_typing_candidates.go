package errors

import (
	"fmt"
)

// ErrMultipleDIDuckTypingCandidates indicates that more than one element was found that fits the given type.
type ErrMultipleDIDuckTypingCandidates struct {
	baseErrorChecker[ErrMultipleDIDuckTypingCandidates, *ErrMultipleDIDuckTypingCandidates]
	Type  any
	Count int
}

func (e ErrMultipleDIDuckTypingCandidates) Error() string {
	givenType := genDebugKeyTypeName(e.Type)

	return fmt.Sprintf("found %d elements that fits the given type", e.Count) + givenType
}
