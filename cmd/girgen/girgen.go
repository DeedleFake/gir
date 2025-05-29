package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"text/template"

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

		"isClass": func(name string) bool {
			return strings.HasSuffix(name, "Class")
		},

		"toCallable": func(info *gi.BaseInfo) *gi.CallableInfo {
			c, _ := gi.TypeCallableInfo.Check(info)
			return c
		},

		"toStruct": func(info *gi.BaseInfo) *gi.StructInfo {
			c, _ := gi.TypeStructInfo.Check(info)
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

type BaseInfoer interface {
	GetName() string
	AsGIBaseInfo() *gi.BaseInfo
}

type Generator struct {
	w io.Writer

	Config        *Config
	Repo          *gi.Repository
	Type, Element BaseInfoer
}

func (gen Generator) Generate(name string, t, element BaseInfoer) (string, error) {
	gen.Type = t
	gen.Element = element
	return "", tmpl.ExecuteTemplate(gen.w, name, gen)
}

func (gen Generator) Package() string {
	return strings.ToLower(util.ParseCPrefix(gen.Repo.GetCPrefix(gen.Config.Namespace)))
}

func (gen Generator) CPrefix() string {
	return util.ParseCPrefix(gen.Repo.GetCPrefix(gen.Config.Namespace))
}

func (gen Generator) CName() (string, error) {
	info := gen.Type.AsGIBaseInfo()

	if info, ok := gi.TypeRegisteredTypeInfo.Check(info); ok {
		return fmt.Sprintf("%v%v", gen.CPrefix(), info.GetName()), nil
	}

	if info, ok := gi.TypeCallableInfo.Check(info); ok {
		return fmt.Sprintf("%v_%v", strings.ToLower(gen.CPrefix()), info.GetName()), nil
	}

	return "", fmt.Errorf("don't know how to get C name of type %q", info.TypeName())
}

func (gen Generator) MethodName(tname, mname string) string {
	return fmt.Sprintf("%v_%v_%v", strings.ToLower(gen.CPrefix()), util.ToSnakeCase(tname), mname)
}

func (gen Generator) Arguments() (string, error) {
	callable := gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()

	args := make([]string, 0, callable.GetNArgs())
	for arg := range callable.GetArgs() {
		args = append(args, fmt.Sprintf("%v %v", arg.GetName(), gen.TypeInfoToGo(arg.GetTypeInfo())))
	}

	return strings.Join(args, ", "), nil
}

func (gen Generator) TypeInfoToGo(info *gi.TypeInfo) string {
	var buf strings.Builder

	tag := info.GetTag()

	buf.WriteString(TypeTagToGo(tag))
	switch tag {
	case gi.TypeTagInterface:
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			buf.WriteString(gen.RegisteredTypeToGo(i))
		}
	}

	return buf.String()
}

func (gen Generator) RegisteredTypeToGo(info *gi.RegisteredTypeInfo) string {
	namespace := info.GetNamespace()
	prefix := strings.ToLower(util.ParseCPrefix(gen.Repo.GetCPrefix(namespace))) + "."
	if namespace == gen.Config.Namespace {
		prefix = ""
	}

	return prefix + info.GetName()
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

	tl, err := r.Require(config.Namespace, config.Version, gi.RepositoryLoadFlagNone)
	if err != nil {
		slog.Error("failed to open typelib", "namespace", config.Namespace, "version", config.Version, "err", err)
		os.Exit(1)
	}
	defer tl.Unref()

	var buf bytes.Buffer
	_, err = Generator{
		w:      &buf,
		Config: config,
		Repo:   r,
	}.Generate("file", nil, nil)
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
