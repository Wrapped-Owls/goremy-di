package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/wrapped-owls/goremy-di/benchmark/cmd/update-readme/render"
	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchparse"
)

func main() {
	readmePath := flag.String("readme", "README.md", "path to benchmark README")
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		log.Fatal("provide at least one benchmark file path")
	}

	readmeRaw, err := os.ReadFile(*readmePath)
	if err != nil {
		log.Fatal("read README: ", err)
	}

	libraries, order, err := benchparse.ParseLibraries(paths)
	if err != nil {
		log.Fatal("parse benchmark files: ", err)
	}

	goos, goarch, cpu := benchparse.PickMeta(libraries)
	section := render.BuildSection(runtime.Version(), goos, goarch, cpu, libraries, order)

	updated, err := replaceSection(string(readmeRaw), section)
	if err != nil {
		log.Fatal("update README markers: ", err)
	}

	if err = os.WriteFile(*readmePath, []byte(updated), 0o644); err != nil {
		log.Fatal("write README: ", err)
	}
}
