package fixtures

type TestContextRepository interface {
	RequestID() string
}

type TestContextRepositoryImpl struct {
	RequestIDValue string
}

func (r TestContextRepositoryImpl) RequestID() string {
	return r.RequestIDValue
}
