package storage

// DepsStorage holds all dependencies
type DepsStorage struct {
	allowOverride  bool
	namedInstances map[string]any
	instances      []any
}

func NewDepsStorage(allowOverride bool) *DepsStorage {
	return &DepsStorage{
		allowOverride:  allowOverride,
		namedInstances: map[string]any{},
		instances:      make([]any, 0, 11),
	}
}

func (s *DepsStorage) AllowOverride(value bool) {
	s.allowOverride = value
}

func (s *DepsStorage) Set(value any) {
	// Change signature to receive the bType here
	s.instances = append(s.instances, value)
}

func (s *DepsStorage) SetNamed(key string, value any) {
	if _, ok := s.namedInstances[key]; ok && !s.allowOverride {
		panic("override not allowed")
	}
	s.namedInstances[key] = value
}

func (s DepsStorage) Get(key string) any {
	result, _ := s.namedInstances[key]
	return result
}

func (s DepsStorage) Binds() []any {
	return s.instances
}
