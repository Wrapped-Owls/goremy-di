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

func (c loggerInjector) BindElem(key remy.BindKey, element any, opts remy.BindOptions) (err error) {
	c.logger.Infof("Injector[BindElem](%v, %v, %#v) - Starting\n", key, element, opts)
	err = c.ij.BindElem(key, element, opts)
	if err != nil {
		c.logger.Errorf("Injector[BindElem]<%v:%#v> - Error: `%v`\n", key, opts, err)
	}
	c.logger.Infof("Injector[BindElem](%v, %v, %#v) - Ending\n", key, element, opts)
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

func (c loggerInjector) RetrieveBind(bindKey remy.BindKey, tag string) (result any, err error) {
	c.logger.Infof("Injector[RetrieveBind](%v, %s) - Starting\n", bindKey, tag)
	if result, err = c.ij.RetrieveBind(bindKey, tag); err != nil {
		c.logger.Errorf("Injector[RetrieveBind]<%v:%s> - Error: `%v`\n", bindKey, tag, err)
		return
	}

	c.logger.Infof("Injector[RetrieveBind](%v, %s) - Found `%+v`\n", bindKey, tag, result)
	return
}

func (c loggerInjector) GetAll(optKey ...string) (result []any, err error) {
	c.logger.Infof("Injector[GetAll](%+v) - Starting\n", optKey)
	result, err = c.ij.GetAll(optKey...)
	if err != nil {
		c.logger.Errorf("Injector[GetAll]<%v> - Error: `%v`\n", optKey, err)
	}

	c.logger.Infof("Injector[GetAll](%+v) - Found `%+v`\n", optKey, result)
	return
}

func (c loggerInjector) ReflectOpts() (opts remy.ReflectionOptions) {
	opts = c.ij.ReflectOpts()
	c.logger.Infof("Injector[ReflectOpts] - Returning `%+v`\n", opts)
	return
}
