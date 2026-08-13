package remychef

import "github.com/wrapped-owls/goremy-di/remy"

// ModuleRegister adds one bind to an injector on behalf of the App that owns it.
// It is the App-aware counterpart of remy.ModuleRegister.
type ModuleRegister func(app *App, inj remy.Injector)

// Module turns registers into a remy.Module bound to this App, so tracked binds
// compose with the plain remy.Module values the application already builds.
func (app *App) Module(registers ...ModuleRegister) remy.Module {
	moduleRegisters := make([]remy.ModuleRegister, 0, len(registers))
	for _, register := range registers {
		moduleRegisters = append(moduleRegisters, app.bind(register))
	}

	return remy.NewModule(moduleRegisters...)
}

// go.mod pins 1.20: taking register as a parameter is what stops every closure
// from capturing the same loop variable
func (app *App) bind(register ModuleRegister) remy.ModuleRegister {
	return func(inj remy.Injector) { register(app, inj) }
}

// Register applies registers straight to the App injector, tracking what they build.
func (app *App) Register(registers ...ModuleRegister) error {
	return app.RegisterModule(app.Module(registers...))
}

// RegisterModule delegates to remy.RegisterModule. A plain remy.Module registers
// binds the App cannot see; Register is the tracked counterpart.
func (app *App) RegisterModule(modules ...remy.Module) error {
	return remy.RegisterModule(app.Injector, modules...)
}
