package main

import (
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func main() {
	var outputDir string
	flag.StringVar(&outputDir, "output", "", "Directory to write generated benchmark files")
	flag.Parse()

	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		slog.Error("failed to create output dir", slog.String("err", err.Error()))
		os.Exit(1)
	}

	for _, lib := range libraries {
		tmpl, err := loadTemplate(lib.TemplateName, template.FuncMap{})
		if err != nil {
			slog.Error(
				"failed to load template",
				slog.String("lib", lib.Name),
				slog.String("err", err.Error()),
			)
			os.Exit(1)
		}

		data, err := buildTemplateData(lib)
		if err != nil {
			slog.Error(
				"failed to build template data",
				slog.String("lib", lib.Name),
				slog.String("err", err.Error()),
			)
			os.Exit(1)
		}

		outputPath := filepath.Join(outputDir, lib.FileName)
		if err = writeTemplateFile(tmpl, outputPath, data); err != nil {
			slog.Error("failed to write", slog.String("err", err.Error()))
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", outputPath)
	}
}

func loadTemplate(templateName string, funcMap template.FuncMap) (*template.Template, error) {
	files := []string{
		"templates/base.go.tmpl",
		filepath.Join("templates", templateName),
	}
	return template.New(templateName).Funcs(funcMap).ParseFS(templateFS, files...)
}

func writeTemplateFile(tmpl *template.Template, outputPath string, data TemplateData) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", outputPath, err)
	}
	defer file.Close()

	if err = tmpl.ExecuteTemplate(file, "file", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
