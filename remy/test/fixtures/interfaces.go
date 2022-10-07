package fixtures

type Language interface {
	Name() string
	Kind() string
}

type GoProgrammingLang struct{}

func (g GoProgrammingLang) Name() string {
	return "Go"
}

func (g GoProgrammingLang) Kind() string {
	return "programming"
}
