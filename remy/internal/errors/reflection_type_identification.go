package errors

// ErrImpossibleIdentifyType indicates that it's impossible to identify the given type.
type ErrImpossibleIdentifyType struct {
	baseErrorChecker[ErrImpossibleIdentifyType, *ErrImpossibleIdentifyType]
	Type any
}

func (e ErrImpossibleIdentifyType) Error() string {
	return "impossible to identify the given type" + genDebugKeyTypeName(e.Type)
}
