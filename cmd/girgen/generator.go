package main

import (
	"fmt"
	"io"
	"strings"
	"unsafe"

	"deedles.dev/gir/g"
	"deedles.dev/gir/gi"
	"deedles.dev/gir/internal/util"
)

type BaseInfoer interface {
	g.TypeInstancer
	GetName() string
	GetNamespace() string
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

func (gen *Generator) PackageFor(namespace string) string {
	pkg := strings.ToLower(util.ParseCPrefix(gen.Repo.GetCPrefix(namespace)))
	if pkg == gen.Package() {
		return ""
	}
	return pkg
}

func (gen *Generator) MethodPrefix() string {
	info := gi.TypeRegisteredTypeInfo.Cast(gen.Type)
	return strings.TrimSuffix(info.GetTypeInitFunctionName(), "_get_type")
}

func (gen *Generator) MethodName(name string) string {
	return fmt.Sprintf("%v_%v", gen.MethodPrefix(), name)
}

func (gen *Generator) InstanceSize() uint {
	if info, ok := gi.TypeObjectInfo.Check(gen.Type); ok {
		size := info.GetGType().Query().InstanceSize()
		if parent := info.GetParent(); parent != nil {
			return size - parent.GetGType().Query().InstanceSize()
		}
		return size - uint(unsafe.Sizeof(*new(g.TypeInstance)))
	}

	if info, ok := gi.TypeRegisteredTypeInfo.Check(gen.Type); ok {
		return info.GetGType().Query().InstanceSize()
	}

	panic(fmt.Errorf("can't get size for %v", gen.Type.AsGTypeInstance().TypeName()))
}

func (gen *Generator) Callable() *gi.CallableInfo {
	return gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()
}

func (gen *Generator) Arguments() *Arguments {
	args := Arguments{Generator: gen}
	args.Load()
	return &args
}

func (gen *Generator) TypeInfoToGo(info *gi.TypeInfo) string {
	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagVoid:
		return "unsafe.Pointer"

	case gi.TypeTagInterface:
		if info.IsPointer() {
			buf.WriteByte('*')
		}
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			pkg := gen.PackageFor(i.GetNamespace())
			if pkg != "" {
				buf.WriteString(pkg)
				buf.WriteByte('.')
			}
			buf.WriteString(i.GetName())
		}

	case gi.TypeTagGtype:
		pkg := gen.PackageFor("GObject")
		if pkg != "" {
			buf.WriteString(pkg)
			buf.WriteByte('.')
		}
		buf.WriteString("Type")

	case gi.TypeTagArray:
		fmt.Fprintf(&buf, "[]%v", gen.TypeInfoToGo(info.GetParamType(0)))

	default:
		buf.WriteString(typeTagsGo[tag])
	}

	return buf.String()
}

func (gen *Generator) TypeInfoToC(info *gi.TypeInfo) string {
	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagVoid:
		if info.IsPointer() {
			buf.WriteString("unsafe.Pointer")
			break
		}
		buf.WriteString(typeTagsC[tag])

	case gi.TypeTagInterface:
		if info.IsPointer() {
			buf.WriteByte('*')
		}
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			buf.WriteString("C.")
			buf.WriteString(CTypeName(i))
		}

	case gi.TypeTagArray:
		fmt.Fprintf(&buf, "*%v", gen.TypeInfoToC(info.GetParamType(0)))

	default:
		buf.WriteString(typeTagsC[tag])
	}

	return buf.String()
}
