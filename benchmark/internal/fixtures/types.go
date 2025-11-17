package fixtures

// Common test types used across all benchmarks
type (
	ServiceA interface {
		DoWork() string
	}

	ServiceB interface {
		Process() string
	}

	ServiceC interface {
		Execute() string
	}

	Repository interface {
		GetData() string
	}

	Config struct {
		Value string
	}

	ServiceAImpl struct {
		Repo Repository
		Cfg  *Config
	}

	ServiceBImpl struct {
		ServiceA ServiceA
		Cfg      *Config
	}

	ServiceCImpl struct {
		ServiceA ServiceA
		ServiceB ServiceB
		Cfg      *Config
	}

	RepositoryImpl struct {
		Cfg *Config
	}
)

// Implementations
func (s *ServiceAImpl) DoWork() string {
	return "ServiceA: " + s.Repo.GetData() + " | " + s.Cfg.Value
}

func (s *ServiceBImpl) Process() string {
	return "ServiceB: " + s.ServiceA.DoWork() + " | " + s.Cfg.Value
}

func (s *ServiceCImpl) Execute() string {
	return "ServiceC: " + s.ServiceA.DoWork() + " + " + s.ServiceB.Process() + " | " + s.Cfg.Value
}

func (r *RepositoryImpl) GetData() string {
	return "Data from " + r.Cfg.Value
}
