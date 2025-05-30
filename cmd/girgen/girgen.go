package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"deedles.dev/gir/g"
	"deedles.dev/gir/gi"
	"deedles.dev/gir/internal/util"
	"golang.org/x/tools/imports"
)

var (
	//go:embed tmpl
	tmplFS embed.FS

	tmpl = template.Must(template.New("root").Funcs(tmplFuncs).ParseFS(tmplFS, "tmpl/*"))

	tmplFuncs = template.FuncMap{
		"toCamelCase": util.ToCamelCase,
		"toSnakeCase": util.ToSnakeCase,

		"parent": func(info *gi.ObjectInfo) parentInfo {
			if parent := info.GetParent(); parent != nil {
				return parent
			}
			return typeInstanceParentInfo{}
		},

		"isGPointerReceiver": func(name string) bool {
			return name == "clear" || name == "ref" || name == "unref"
		},

		"isClass": func(name string) bool {
			return strings.HasSuffix(name, "Class")
		},

		"toCallable": func(info g.TypeInstancer) *gi.CallableInfo {
			c, _ := gi.TypeCallableInfo.Check(info)
			return c
		},

		"toStruct": func(info g.TypeInstancer) *gi.StructInfo {
			c, _ := gi.TypeStructInfo.Check(info)
			return c
		},

		"toObject": func(info g.TypeInstancer) *gi.ObjectInfo {
			c, _ := gi.TypeObjectInfo.Check(info)
			return c
		},

		"toConstant": func(info g.TypeInstancer) *gi.ConstantInfo {
			c, _ := gi.TypeConstantInfo.Check(info)
			return c
		},

		"toEnum": func(info g.TypeInstancer) *gi.EnumInfo {
			c, _ := gi.TypeEnumInfo.Check(info)
			return c
		},

		"toFlags": func(info g.TypeInstancer) *gi.FlagsInfo {
			c, _ := gi.TypeFlagsInfo.Check(info)
			return c
		},
	}
)

func readConfig(path string) *Config {
	slog := slog.With("path", path)

	file, err := os.Open(path)
	if err != nil {
		slog.Error("failed to open config file", "err", err)
		os.Exit(1)
	}
	defer file.Close()

	config, err := ParseConfig(file)
	if err != nil {
		slog.Error("failed to parse config", "err", err)
		os.Exit(1)
	}

	return config
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: girgen -o output.go config.gen")
	}
	output := flag.String("o", "", "output filename")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	config := readConfig(flag.Arg(0))

	r := gi.RepositoryNew()

	tl, err := r.Require(config.Namespace, config.Version, gi.RepositoryLoadFlagsNone)
	if err != nil {
		slog.Error("failed to open typelib", "namespace", config.Namespace, "version", config.Version, "err", err)
		os.Exit(1)
	}
	defer tl.Unref()

	var buf bytes.Buffer
	err = Generate(&buf, config, r)
	if err != nil {
		slog.Error("failed to execute template", "err", err)
		os.Exit(1)
	}

	formatted, err := imports.Process(*output, buf.Bytes(), nil)
	if err != nil {
		os.Stdout.Write(buf.Bytes())
		slog.Error("failed to format output", "err", err)
		os.Exit(1)
	}

	if *output == "" {
		os.Stdout.Write(formatted)
		return
	}

	err = os.WriteFile(*output, formatted, 0644)
	if err != nil {
		slog.Error("failed to write output", "err", err)
		os.Exit(1)
	}
}
