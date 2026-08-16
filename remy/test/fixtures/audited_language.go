package fixtures

// AuditLog collects the names an AuditedLanguage observed
type AuditLog struct {
	Lookups []string
}

func (a *AuditLog) Record(name string) {
	a.Lookups = append(a.Lookups, name)
}

// AuditedLanguage decorates a Language so every name lookup reaches the audit log
type AuditedLanguage struct {
	Inner Language
	Audit *AuditLog
}

func (l AuditedLanguage) Name() string {
	name := l.Inner.Name()
	l.Audit.Record(name)
	return name
}

func (l AuditedLanguage) Kind() string {
	return l.Inner.Kind()
}
