package fixtures

// Constructor functions for all DI libraries
func NewConfig() *Config {
	return &Config{Value: "test-config"}
}

func NewRepository(cfg *Config) Repository {
	return &RepositoryImpl{Cfg: cfg}
}

func NewServiceA(repo Repository, cfg *Config) ServiceA {
	return &ServiceAImpl{Repo: repo, Cfg: cfg}
}

func NewServiceB(serviceA ServiceA, cfg *Config) ServiceB {
	return &ServiceBImpl{ServiceA: serviceA, Cfg: cfg}
}

func NewServiceC(serviceA ServiceA, serviceB ServiceB, cfg *Config) ServiceC {
	return &ServiceCImpl{ServiceA: serviceA, ServiceB: serviceB, Cfg: cfg}
}

// Constructor functions that return concrete types (for wire)
func NewRepositoryImpl(cfg *Config) *RepositoryImpl {
	return &RepositoryImpl{Cfg: cfg}
}

func NewServiceAImpl(repo Repository, cfg *Config) *ServiceAImpl {
	return &ServiceAImpl{Repo: repo, Cfg: cfg}
}

func NewServiceBImpl(serviceA ServiceA, cfg *Config) *ServiceBImpl {
	return &ServiceBImpl{ServiceA: serviceA, Cfg: cfg}
}

func NewServiceCImpl(serviceA ServiceA, serviceB ServiceB, cfg *Config) *ServiceCImpl {
	return &ServiceCImpl{ServiceA: serviceA, ServiceB: serviceB, Cfg: cfg}
}
