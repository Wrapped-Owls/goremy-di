package fixtures

type ServiceAImpl struct {
	Repo Repository
	Cfg  *Config
}

func (s ServiceAImpl) DoWork() string {
	return "ServiceA: " + s.Repo.GetData() + " | " + s.Cfg.Value
}

func NewServiceA(repo Repository, cfg *Config) ServiceA {
	return ServiceAImpl{Repo: repo, Cfg: cfg}
}
