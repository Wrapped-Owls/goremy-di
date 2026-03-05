package render

type benchBest struct {
	opsMax   int64
	nsMin    float64
	bMin     float64
	allocMin float64
}

type benchSection struct {
	Title string
	Rows  []benchRow
}

type benchRow struct {
	Library string
	Ops     string
	NsOp    string
	BOp     string
	Allocs  string
}
