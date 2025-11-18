package fixtures

type ServiceBImpl struct {
	ServiceA ServiceA
	Cfg      *Config
}

func (s ServiceBImpl) Process() string {
	return "ServiceB: " + s.ServiceA.DoWork() + " | " + s.Cfg.Value
}

func NewServiceB(serviceA ServiceA, cfg *Config) ServiceB {
	return ServiceBImpl{ServiceA: serviceA, Cfg: cfg}
}
