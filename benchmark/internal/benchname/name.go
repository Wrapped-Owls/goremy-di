package benchname

import "strings"

func HumanizeBenchmark(name string) string {
	replacer := strings.NewReplacer("_", " ", "/", " ")
	name = replacer.Replace(name)

	parts := strings.Fields(name)
	for i := range parts {
		parts[i] = splitCamel(parts[i])
	}

	return strings.Join(parts, " ")
}

func HumanizeLibrary(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Unknown"
	}

	replacer := strings.NewReplacer("_", " ", "-", " ")
	parts := strings.Fields(replacer.Replace(name))
	for i, part := range parts {
		part = splitCamel(part)
		if part == "" {
			continue
		}
		frags := strings.Fields(strings.ToLower(part))
		for j := range frags {
			frags[j] = strings.ToUpper(frags[j][:1]) + frags[j][1:]
		}
		parts[i] = strings.Join(frags, " ")
	}
	return strings.Join(parts, " ")
}

func splitCamel(word string) string {
	if word == "" {
		return word
	}

	var out strings.Builder
	for i, r := range word {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(word[i-1])
			if prev >= 'a' && prev <= 'z' {
				out.WriteByte(' ')
			}
		}
		out.WriteRune(r)
	}
	return out.String()
}
