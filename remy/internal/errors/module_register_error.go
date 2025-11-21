package errors

import "fmt"

type ErrModuleRegisterErrors struct {
	baseErrorChecker[ErrModuleRegisterErrors, *ErrModuleRegisterErrors]
	MainMessage string
	ErrList     []error
}

func NewErrModuleRegisterErrors(mainMessage string, errList ...error) ErrModuleRegisterErrors {
	if mainMessage == "" {
		mainMessage = "module register error"
	}
	return ErrModuleRegisterErrors{MainMessage: mainMessage, ErrList: errList}
}

func (e ErrModuleRegisterErrors) Unwrap() []error {
	return e.ErrList
}

func (e ErrModuleRegisterErrors) Error() string {
	return fmt.Sprintf("%s, failed to register %d modules", e.MainMessage, len(e.ErrList))
}
