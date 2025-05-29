package main

import (
	"fmt"
	"iter"
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

	var remove []int
	for i, info := range callable.GetArgs() {
		arg := Argument{
			Generator: args.Generator,
			Index:     i,
			Info:      info,
		}
		args.Args = append(args.Args, &arg)
		remove = append(remove, arg.SubArgs()...)
	}

	slices.Sort(remove)
	for _, i := range slices.Backward(remove) {
		args.Args = slices.Delete(args.Args, i, i)
	}
}

func (args *Arguments) GoInput() string {
	return util.JoinPairs(argsGoNameType(args.goInput()), " ", ", ")
}

func (args *Arguments) GoOutput() string {
	return util.JoinPairs(argsGoNameType(args.goOutput()), " ", ", ")
}

func (args *Arguments) CInput() string {
	return util.JoinSeq(xiter.Flatten(argsCNames(args.cInput())), ", ")
}

func (args *Arguments) ConvertToC() string {
	callable := args.Callable()

	var buf strings.Builder
	for i, arg := range callable.GetArgs() {
		ti := arg.GetTypeInfo()
		switch tag := ti.GetTag(); tag {
		case gi.TypeTagUtf, gi.TypeTagFilename:
			fmt.Fprintf(&buf, "arg%v := C.CString(%v)\ndefer C.free(unsafe.Pointer(arg%v))\n", i, arg.GetName(), i)
		default:
			fmt.Fprintf(&buf, "arg%v := (%v)(%v)", i, args.TypeInfoToC(ti), arg.GetName())
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

type Argument struct {
	*Generator

	Index uint
	Info  *gi.ArgInfo
}

func (arg *Argument) SubArgs() []int {
	// TODO
	return nil
}

func (arg *Argument) IsInput() bool {
	return !arg.Info.IsReturnValue() && arg.Info.GetDirection() != gi.DirectionOut
}

func (arg *Argument) IsCInput() bool {
	return !arg.Info.IsReturnValue()
}

func (arg *Argument) GoName() string {
	return arg.Info.GetName()
}

func (arg *Argument) GoType() string {
	info := arg.Info.GetTypeInfo()

	var buf strings.Builder
	switch tag := info.GetTag(); tag {
	case gi.TypeTagInterface:
		if info.IsPointer() {
			buf.WriteByte('*')
		}
		i := info.GetInterface()
		if i, ok := gi.TypeRegisteredTypeInfo.Check(i); ok {
			buf.WriteString(arg.RegisteredTypeToGo(i))
		}
	default:
		buf.WriteString(TypeTagToGo(tag))
	}

	return buf.String()
}

func (arg *Argument) CNames() []string {
	return []string{fmt.Sprintf("arg%v", arg.Index)}
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

func argsGoNameType(seq iter.Seq[*Argument]) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for arg := range seq {
			if !yield(arg.GoName(), arg.GoType()) {
				return
			}
		}
	}
}

func argsCNames(seq iter.Seq[*Argument]) iter.Seq[iter.Seq[string]] {
	return func(yield func(iter.Seq[string]) bool) {
		for arg := range seq {
			if !yield(slices.Values(arg.CNames())) {
				return
			}
		}
	}
}
