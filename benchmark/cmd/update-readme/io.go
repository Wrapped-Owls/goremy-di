package main

import (
	"fmt"
	"strings"

	"github.com/wrapped-owls/goremy-di/benchmark/cmd/update-readme/render"
)

func replaceSection(readme, section string) (string, error) {
	start := strings.Index(readme, render.MarkerStart)
	end := strings.Index(readme, render.MarkerEnd)

	if start == -1 || end == -1 {
		trimmed := strings.TrimRight(readme, "\n")
		return trimmed + "\n\n" + section + "\n", nil
	}
	if end < start {
		return "", fmt.Errorf("invalid marker order")
	}

	end += len(render.MarkerEnd)
	return readme[:start] + section + readme[end:], nil
}
