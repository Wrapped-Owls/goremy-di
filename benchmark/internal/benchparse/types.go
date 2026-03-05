package benchparse

type Metrics struct {
	Ops    int64
	NsOp   float64
	BOp    float64
	Allocs float64
}

type LibraryData struct {
	Name    string
	Goos    string
	Goarch  string
	CPU     string
	Benches map[string]Metrics
}
