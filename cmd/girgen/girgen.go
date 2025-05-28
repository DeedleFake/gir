package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"text/template"

	"deedles.dev/gir/gi"
)

var (
	//go:embed tmpl
	tmplFS embed.FS

	tmpl = template.Must(template.ParseFS(tmplFS, "tmpl/*"))
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: girgen < config > output.go")
	}
	flag.Parse()

	config, err := ParseConfig(os.Stdin)
	if err != nil {
		slog.Error("failed to parse config", "err", err)
		os.Exit(1)
	}

	r := gi.RepositoryNew()

	tl, err := r.Require(config.Namespace, config.Version, gi.RepositoryLoadFlagNone)
	if err != nil {
		slog.Error("failed to open typelib", "namespace", config.Namespace, "version", config.Version, "err", err)
		os.Exit(1)
	}
	defer tl.Unref()

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "file", map[string]any{
		"Config": config,
	})
	if err != nil {
		slog.Error("failed to execute template", "err", err)
		os.Exit(1)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		os.Stdout.Write(buf.Bytes())
		slog.Error("failed to format output", "err", err)
		os.Exit(1)
	}

	_, err = os.Stdout.Write(formatted)
	if err != nil {
		slog.Error("failed to write output", "err", err)
		os.Exit(1)
	}
}
