package benchformat

import (
	"math"
	"strconv"
	"strings"
)

func Number(v float64) string {
	rounded := math.Round(v)
	if math.Abs(v-rounded) < 1e-9 {
		return Int(int64(rounded))
	}

	s := strconv.FormatFloat(v, 'f', 3, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func Int(v int64) string {
	s := strconv.FormatInt(v, 10)
	if len(s) <= 3 {
		return s
	}

	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	prefix := len(s) % 3
	if prefix == 0 {
		prefix = 3
	}

	var out strings.Builder
	if neg {
		out.WriteByte('-')
	}
	out.WriteString(s[:prefix])
	for i := prefix; i < len(s); i += 3 {
		out.WriteByte(',')
		out.WriteString(s[i : i+3])
	}
	return out.String()
}
