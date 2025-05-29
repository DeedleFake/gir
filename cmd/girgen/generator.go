package main

import (
	"fmt"
	"io"
	"strings"

	"deedles.dev/gir/gi"
	"deedles.dev/gir/internal/util"
)

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

func Generate(w io.Writer, config *Config, r *gi.Repository) error {
	_, err := Generator{
		w:      w,
		Config: config,
		Repo:   r,
	}.Generate("file", nil, nil)
	return err
}

func (gen Generator) Generate(name string, t, element BaseInfoer) (string, error) {
	gen.Type = t
	gen.Element = element
	return "", tmpl.ExecuteTemplate(gen.w, name, &gen)
}

func (gen *Generator) Package() string {
	return strings.ToLower(gen.CPrefix())
}

func (gen *Generator) CPrefix() string {
	return util.ParseCPrefix(gen.Repo.GetCPrefix(gen.Config.Namespace))
}

func (gen *Generator) CName() (string, error) {
	info := gen.Type.AsGIBaseInfo()

	if info, ok := gi.TypeRegisteredTypeInfo.Check(info); ok {
		return fmt.Sprintf("%v%v", gen.CPrefix(), info.GetName()), nil
	}

	if info, ok := gi.TypeCallableInfo.Check(info); ok {
		return fmt.Sprintf("%v_%v", strings.ToLower(gen.CPrefix()), info.GetName()), nil
	}

	return "", fmt.Errorf("don't know how to get C name of type %q", info.TypeName())
}

func (gen *Generator) MethodName(tname, mname string) string {
	return util.MethodName(gen.CPrefix(), tname, mname)
}

func (gen *Generator) Callable() *gi.CallableInfo {
	return gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()
}

func (gen *Generator) Arguments() *Arguments {
	args := Arguments{Generator: gen}
	args.Load()
	return &args
}
