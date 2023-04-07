package utils

import (
	"github.com/wrapped-owls/goremy-di/remy"
)

type loggerInjector struct {
	ij     remy.Injector
	logger Logger
}

func NewLoggerInjector(logger Logger, configs ...remy.Config) *loggerInjector {
	return &loggerInjector{ij: remy.NewInjector(configs...), logger: logger}
}

func (c loggerInjector) Bind(key remy.BindKey, element any) (err error) {
	c.logger.Infof("Injector[Bind](%v, %v) - Starting\n", key, element)
	err = c.ij.Bind(key, element)
	if err != nil {
		c.logger.Errorf("Injector[Bind] - Error: `%v`\n", err)
	}
	c.logger.Infof("Injector[Bind](%v, %v) - Ending\n", key, element)
	return
}

func (c loggerInjector) BindNamed(key remy.BindKey, name string, element any) (err error) {
	c.logger.Infof("Injector[Bind](%v, %s, %v) - Starting\n", key, name, element)
	err = c.ij.BindNamed(key, name, element)
	if err != nil {
		c.logger.Errorf("Injector[BindNamed] - Error: `%v`\n", err)
	}
	c.logger.Infof("Injector[Bind](%v, %s, %v) - Ending\n", key, name, element)
	return
}

func (c loggerInjector) SubInjector(allowOverrides ...bool) remy.Injector {
	var shouldOverride bool
	if len(allowOverrides) > 0 {
		shouldOverride = allowOverrides[0]
	}

	var (
		parentConfig = c.ij.ReflectOpts()
		config       = remy.Config{
			ParentInjector:     c,
			CanOverride:        shouldOverride,
			GenerifyInterfaces: parentConfig.GenerifyInterface,
			UseReflectionType:  parentConfig.UseReflectionType,
		}
	)

	c.logger.Infof("Creating SubInjector with: `%+v`\n", config)
	inj := NewLoggerInjector(c.logger, config)
	return inj
}

func (c loggerInjector) WrapRetriever() remy.Injector {
	c.logger.Info("Injector[WrapRetriever] - Returning nil")
	return nil
}

func (c loggerInjector) GetNamed(key remy.BindKey, name string) (result any, err error) {
	c.logger.Infof("Injector[GetNamed](%v, %s) - Starting\n", key, name)
	result, err = c.ij.GetNamed(key, name)
	if err != nil {
		c.logger.Errorf("Injector[GetNamed] - Error: `%v`\n", err)
	}

	c.logger.Infof("Injector[GetNamed](%v, %s) - Found `%+v`\n", key, name, result)
	return
}

func (c loggerInjector) Get(key remy.BindKey) (result any, err error) {
	c.logger.Infof("Injector[Get](%v) - Starting\n", key)
	result, err = c.ij.Get(key)
	if err != nil {
		c.logger.Errorf("Injector[Get] - Error: `%v`\n", err)
	}

	c.logger.Infof("Injector[Get](%v) - Found `%+v`\n", key, result)
	return
}

func (c loggerInjector) GetAll(optKey ...string) (result []any, err error) {
	c.logger.Infof("Injector[GetAll](%+v) - Starting\n", optKey)
	result, err = c.ij.GetAll(optKey...)
	if err != nil {
		c.logger.Errorf("Injector[GetAll] - Error: `%v`\n", err)
	}

	c.logger.Infof("Injector[GetAll](%+v) - Found `%+v`\n", optKey, result)
	return
}

func (c loggerInjector) ReflectOpts() (opts remy.ReflectionOptions) {
	opts = c.ij.ReflectOpts()
	c.logger.Infof("Injector[ReflectOpts] - Returning `%+v`\n", opts)
	return
}
