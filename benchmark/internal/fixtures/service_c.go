package fixtures

type ServiceCImpl struct {
	ServiceA ServiceA
	ServiceB ServiceB
	Cfg      *Config
}

func (s ServiceCImpl) Execute() string {
	return "ServiceC: " + s.ServiceA.DoWork() + " + " + s.ServiceB.Process() + " | " + s.Cfg.Value
}

func NewServiceC(serviceA ServiceA, serviceB ServiceB, cfg *Config) ServiceC {
	return ServiceCImpl{ServiceA: serviceA, ServiceB: serviceB, Cfg: cfg}
}
