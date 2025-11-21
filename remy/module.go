package remy

import remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"

// Module represents a dependency-injection module.
// A module encapsulates a set of provider registrations and is expected to register all of its bindings when
// Register is called with the given Injector.
type Module interface {
	Register(injector Injector)
}

// RegisterModuleFunc registers one or more registration functions against the provided Injector.
// Each function receives the Injector instance and is expected to perform its own provider registrations.
func RegisterModuleFunc(inj Injector, modules ...func(Injector)) (err error) {
	if inj == nil {
		return remyErrs.NewErrModuleRegisterErrors("injector is required")
	}
	defer recoverInjectorPanic(&err)

	for _, module := range modules {
		if module != nil {
			module(inj)
		}
	}

	return nil
}

// RegisterModule registers one or more Module instances using the provided Injector.
// Each Module's Register method is adapted and delegated to RegisterModuleFunc.
func RegisterModule(inj Injector, modules ...Module) (err error) {
	inj = mustInjector(inj)
	var errList []error
	for _, module := range modules {
		if err = RegisterModuleFunc(inj, module.Register); err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		return remyErrs.NewErrModuleRegisterErrors("", errList...)
	}

	return err
}

// ModuleRegister is a helper function signature that can be passed to NewModule.
// It receives an Injector and performs one or more registrations on it.
type ModuleRegister func(Injector)

// simpleModule is a lightweight Module implementation that aggregates multiple registrations.
type simpleModule struct{ regs []ModuleRegister }

func (m *simpleModule) Register(injector Injector) {
	for _, r := range m.regs {
		if r != nil {
			r(injector)
		}
	}

	m.regs = nil // Allow Go runtime to delete temporary register functions
}

// NewModule creates a Module from a list of ModuleRegister functions.
//
// Usage example:
//
//	userModule := remy.NewModule(
//	    remy.WithConstructor(remy.Factory[UserService], NewUserService),
//	    remy.WithInstance("static value"),
//	    remy.WithBind(remy.Factory(func(d remy.DependencyRetriever) (Svc, error) { ... })),
//	)
//	remy.RegisterModule(inj, userModule)
func NewModule(registers ...ModuleRegister) Module {
	// copy to avoid external modification of internal slice
	return &simpleModule{regs: registers}
}
