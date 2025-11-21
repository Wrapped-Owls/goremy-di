package remy

type Module interface {
	Register(injector Injector)
}

func RegisterModule(inj Injector, modules ...Module) (err error) {
	defer recoverInjectorPanic(&err)
	inj = mustInjector(inj)
	for _, module := range modules {
		module.Register(inj)
	}

	return nil
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
