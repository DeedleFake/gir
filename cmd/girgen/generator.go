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
	return "", tmpl.ExecuteTemplate(gen.w, name, gen)
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

func (gen *Generator) Arguments() *Arguments {
	callable := gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()

	args := Arguments{Gen: gen}
	for i, arg := range callable.GetArgs() {
		args.Args = append(args.Args, &Argument{
			Index: i,
			Info:  arg,
		})
	}

	return &args
}

func (gen *Generator) TypeInfoToGo(info *gi.TypeInfo) string {
	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagInterface:
		if info.IsPointer() {
			buf.WriteByte('*')
		}
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			buf.WriteString(gen.RegisteredTypeToGo(i))
		}
	default:
		buf.WriteString(TypeTagToGo(tag))
	}

	return buf.String()
}

func (gen *Generator) RegisteredTypeToGo(info *gi.RegisteredTypeInfo) string {
	localPrefix := strings.ToLower(gen.CPrefix()) + "."
	typePrefix := strings.ToLower(util.ParseCPrefix(gen.Repo.GetCPrefix(info.GetNamespace()))) + "."
	if localPrefix == typePrefix {
		typePrefix = ""
	}

	return typePrefix + info.GetName()
}

func (gen *Generator) TypeInfoToC(info *gi.TypeInfo) string {
	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagInterface:
		if info.IsPointer() {
			buf.WriteByte('*')
		}
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			buf.WriteString("C.")
			buf.WriteString(i.GetTypeName())
		}

	default:
		buf.WriteString(TypeTagToC(tag))
	}

	return buf.String()
}

func (gen *Generator) ConvertArguments() string {
	callable := gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()

	var buf strings.Builder
	for i, arg := range callable.GetArgs() {
		ti := arg.GetTypeInfo()
		switch tag := ti.GetTag(); tag {
		case gi.TypeTagUtf, gi.TypeTagFilename:
			fmt.Fprintf(&buf, "arg%v := C.CString(%v)\ndefer C.free(unsafe.Pointer(arg%v))\n", i, arg.GetName(), i)
		default:
			fmt.Fprintf(&buf, "arg%v := (%v)(%v)", i, gen.TypeInfoToC(ti), arg.GetName())
		}
	}

	return buf.String()
}

func (gen *Generator) Call(receiver string) string {
	callable := gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()

	args := make([]string, 0, callable.GetNArgs())
	for i := range callable.GetNArgs() {
		var address string
		if callable.GetArg(i).GetDirection() != gi.DirectionIn {
			address = "&"
		}
		args = append(args, fmt.Sprintf("%varg%v", address, i))
	}

	return fmt.Sprintf(
		"C.%v(%v.c(), %v)",
		gen.MethodName(gen.Type.GetName(), gen.Element.GetName()),
		receiver,
		strings.Join(args, ", "),
	)
}
