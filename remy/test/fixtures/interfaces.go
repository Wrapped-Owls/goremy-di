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

type CountryLanguage struct {
	Language string
}

func (l CountryLanguage) Name() string {
	return l.Language
}

func (l CountryLanguage) Kind() string {
	return "idiom"
}

// PointerLanguage satisfies Language only through a pointer receiver, which is what
// makes the duck typing scan fall back to PointerValue
type PointerLanguage struct {
	Language string
}

func (l *PointerLanguage) Name() string {
	return l.Language
}

func (l *PointerLanguage) Kind() string {
	return "pointer"
}
