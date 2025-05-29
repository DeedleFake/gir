package main

import (
	"fmt"
	"iter"
	"runtime"
	"slices"
	"strings"

	"deedles.dev/gir/gi"
	"deedles.dev/gir/internal/util"
	"deedles.dev/xiter"
)

type Arguments struct {
	*Generator
	Args []*Argument
}

func (args *Arguments) Load() {
	callable := args.Callable()

	var hide []int
	for i, info := range callable.GetArgs() {
		arg := Argument{
			Generator: args.Generator,
			Index:     i,
			Info:      info,
		}
		runtime.AddCleanup(&arg, (*gi.BaseInfo).Unref, info.AsGIBaseInfo())

		args.Args = append(args.Args, &arg)
		hide = append(hide, arg.Obscured()...)
	}

	for _, i := range hide {
		args.Args[i].Hidden = true
	}
}

func (args *Arguments) GoInput() string {
	return util.JoinSeq(argsGoInput(args.goInput()), ", ")
}

func (args *Arguments) GoOutput() string {
	return util.JoinSeq(argsGoInput(args.goOutput()), ", ")
}

func (args *Arguments) CInput() string {
	return util.JoinSeq(argsCInput(args.cInput()), ", ")
}

func (args *Arguments) COutput() string {
	arg := args.cOutput()
	if arg == nil {
		return ""
	}
	return arg.CName()
}

func (args *Arguments) ConvertToC() string {
	var buf strings.Builder
	for arg := range args.cInput() {
		ti := arg.Info.GetTypeInfo()
		switch tag := ti.GetTag(); tag {
		case gi.TypeTagUtf:
			fmt.Fprintf(&buf, "%v := C.CString(%v)\ndefer C.free(unsafe.Pointer(%v))\n", arg.CName(), arg.GoName(), arg.CName())

		default:
			fmt.Fprintf(&buf, "%v := (%v)(%v)", arg.CName(), arg.CType(), arg.CName())
		}
	}

	return buf.String()
}

func (args *Arguments) goInput() iter.Seq[*Argument] {
	return xiter.Filter(slices.Values(args.Args), (*Argument).IsInput)
}

func (args *Arguments) goOutput() iter.Seq[*Argument] {
	return xiter.Filter(slices.Values(args.Args), util.Not((*Argument).IsInput))
}

func (args *Arguments) cInput() iter.Seq[*Argument] {
	return xiter.Filter(slices.Values(args.Args), (*Argument).IsCInput)
}

func (args *Arguments) cOutput() *Argument {
	f := slices.Collect(xiter.Filter(slices.Values(args.Args), func(arg *Argument) bool { return arg.Info.IsReturnValue() }))
	if len(f) == 0 {
		return nil
	}
	return f[0]
}

type Argument struct {
	*Generator

	Index  uint
	Info   *gi.ArgInfo
	Hidden bool
}

func (arg *Argument) Obscured() []int {
	// TODO
	return nil
}

func (arg *Argument) IsInput() bool {
	return !arg.Info.IsReturnValue() && arg.Info.GetDirection() != gi.DirectionOut
}

func (arg *Argument) IsCInput() bool {
	return !arg.Info.IsReturnValue()
}

func (arg *Argument) GoInput() string {
	return fmt.Sprintf("%v %v", arg.GoName(), arg.GoType())
}

func (arg *Argument) GoName() string {
	return arg.Info.GetName()
}

func (arg *Argument) GoType() string {
	info := arg.Info.GetTypeInfo()

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
			buf.WriteString(arg.RegisteredTypeToGo(i))
		}

	default:
		buf.WriteString(typeTagsGo[tag])
	}

	return buf.String()
}

func (arg *Argument) CInput() string {
	var address string
	if arg.Info.GetDirection() != gi.DirectionIn {
		address = "&"
	}
	return fmt.Sprintf("%v%v", address, arg.CName())
}

func (arg *Argument) CType() string {
	info := arg.Info.GetTypeInfo()

	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagVoid:
		if info.IsPointer() {
			buf.WriteString("*C.void")
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
			buf.WriteString(i.GetTypeName())
		}

	default:
		buf.WriteString(typeTagsC[tag])
	}

	return buf.String()
}

func (arg *Argument) CName() string {
	return fmt.Sprintf("arg%v", arg.Index)
}

func (gen *Generator) RegisteredTypeToGo(info *gi.RegisteredTypeInfo) string {
	localPrefix := strings.ToLower(gen.CPrefix()) + "."
	typePrefix := strings.ToLower(util.ParseCPrefix(gen.Repo.GetCPrefix(info.GetNamespace()))) + "."
	if localPrefix == typePrefix {
		typePrefix = ""
	}

	return typePrefix + info.GetName()
}

func argsGoInput(seq iter.Seq[*Argument]) iter.Seq[string] {
	return xiter.Map(seq, (*Argument).GoInput)
}

func argsCInput(seq iter.Seq[*Argument]) iter.Seq[string] {
	return xiter.Map(seq, (*Argument).CInput)
}
