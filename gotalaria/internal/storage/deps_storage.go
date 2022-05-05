package storage

// DepsStorage holds all dependencies
type DepsStorage struct {
	namedBinds map[string]any
	binds      []any
}

func NewDepsStorage() *DepsStorage {
	return &DepsStorage{
		namedBinds: map[string]any{},
		binds:      make([]any, 0, 11),
	}
}

func (s *DepsStorage) Set(value any) {
	s.binds = append(s.binds, value)
}

func (s *DepsStorage) SetNamed(key string, value any) {
	s.namedBinds[key] = value
}

func (s DepsStorage) Get(key string) any {
	result, _ := s.namedBinds[key]
	return result
}

func (s DepsStorage) Binds() []any {
	return s.binds
}
