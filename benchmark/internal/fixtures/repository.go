package fixtures

type RepositoryImpl struct {
	Cfg *Config
}

func (r RepositoryImpl) GetData() string {
	return "Data from " + r.Cfg.Value
}

func NewRepository(cfg *Config) Repository {
	return RepositoryImpl{Cfg: cfg}
}
