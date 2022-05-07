package storage

// DepsStorage holds all dependencies
type DepsStorage struct {
	allowOverride bool
	namedBinds    map[string]any
	binds         []any
}

func NewDepsStorage(allowOverride bool) *DepsStorage {
	return &DepsStorage{
		allowOverride: allowOverride,
		namedBinds:    map[string]any{},
		binds:         make([]any, 0, 11),
	}
}

func (s *DepsStorage) AllowOverride(value bool) {
	s.allowOverride = value
}

func (s *DepsStorage) Set(value any) {
	s.binds = append(s.binds, value)
}

func (s *DepsStorage) SetNamed(key string, value any) {
	if _, ok := s.namedBinds[key]; ok && !s.allowOverride {
		panic("override not allowed")
	}
	s.namedBinds[key] = value
}

func (s DepsStorage) Get(key string) any {
	result, _ := s.namedBinds[key]
	return result
}

func (s DepsStorage) Binds() []any {
	return s.binds
}
