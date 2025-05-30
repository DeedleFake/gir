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

func (args *Arguments) ConvertToGo() string {
	return util.JoinSeq(xiter.Map(args.goOutput(), (*Argument).ConvertToGo), "\n")
}

func (args *Arguments) ConvertToC() string {
	return util.JoinSeq(xiter.Map(args.cInput(), (*Argument).ConvertToC), "\n")
}

func (args *Arguments) goInput() iter.Seq[*Argument] {
	return xiter.Filter(slices.Values(args.Args), (*Argument).IsInput)
}

func (args *Arguments) goOutput() iter.Seq[*Argument] {
	output := xiter.Filter(slices.Values(args.Args), util.Not((*Argument).IsInput))
	if r := args.cOutput(); r != nil {
		output = xiter.Concat(output, xiter.Of(r))
	}
	if err := args.err(); err != nil {
		output = xiter.Concat(output, xiter.Of(err))
	}
	return output
}

func (args *Arguments) cInput() iter.Seq[*Argument] {
	input := xiter.Filter(slices.Values(args.Args), (*Argument).IsCInput)
	if err := args.err(); err != nil {
		input = xiter.Concat(input, xiter.Of(err))
	}
	return input
}

func (args *Arguments) cOutput() *Argument {
	r := args.Callable().GetReturnType()
	if r.GetTag() == gi.TypeTagVoid {
		r.Unref()
		return nil
	}
	return &Argument{Generator: args.Generator, Return: r}
}

func (args *Arguments) err() *Argument {
	r := args.Callable().CanThrowGerror()
	if !r {
		return nil
	}
	return &Argument{Generator: args.Generator, Error: true}
}

type Argument struct {
	*Generator

	Index  uint
	Info   *gi.ArgInfo
	Return *gi.TypeInfo
	Hidden bool
	Error  bool
}

func (arg *Argument) Obscured() []int {
	// TODO
	return nil
}

func (arg *Argument) TypeInfo() *gi.TypeInfo {
	if arg.Return != nil {
		return arg.Return
	}
	return arg.Info.GetTypeInfo()
}

func (arg *Argument) IsInput() bool {
	return arg.Return == nil && !arg.Error && !arg.Info.IsReturnValue() && arg.Info.GetDirection() != gi.DirectionOut
}

func (arg *Argument) IsCInput() bool {
	return arg.Return == nil
}

func (arg *Argument) GoInput() string {
	return fmt.Sprintf("%v %v", arg.GoName(), arg.GoType())
}

func (arg *Argument) GoName() string {
	if arg.Return != nil {
		return "r"
	}
	if arg.Error {
		return "err"
	}
	return arg.Info.GetName()
}

func (arg *Argument) GoType() string {
	if arg.Error {
		return "error"
	}

	info := arg.TypeInfo()

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
	if arg.Error || arg.Info.GetDirection() != gi.DirectionIn {
		address = "&"
	}
	return fmt.Sprintf("%v%v", address, arg.CName())
}

func (arg *Argument) CName() string {
	if arg.Return != nil {
		return "cr"
	}
	if arg.Error {
		return "gerr"
	}
	return fmt.Sprintf("arg%v", arg.Index)
}

func (arg *Argument) CType() string {
	if arg.Error {
		return "*C.GError"
	}

	info := arg.TypeInfo()

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
			buf.WriteString(i.GetTypeName())
		}

	default:
		buf.WriteString(typeTagsC[tag])
	}

	return buf.String()
}

func (arg *Argument) ConvertToGo() string {
	if arg.Error {
		// TODO: Move the error into Go so that it can be freed by the
		// garbage collector. This is an exception to the manual memory
		// management because handling it when it's hidden behind an error
		// interface would be too gosh darn annoying otherwise.
		return fmt.Sprintf("%v = (*g.Error)(unsafe.Pointer(%v))", arg.GoName(), arg.CName())
	}

	ti := arg.TypeInfo()
	switch tag := ti.GetTag(); tag {
	case gi.TypeTagBoolean:
		return fmt.Sprintf("%v = %v != 0", arg.GoName(), arg.CName())

	case gi.TypeTagUtf:
		return fmt.Sprintf("%v = C.GoString(%v)", arg.GoName(), arg.CName())

	case gi.TypeTagInterface:
		return fmt.Sprintf("%v = (%v)(unsafe.Pointer(%v))", arg.GoName(), arg.GoType(), arg.CName())

	default:
		return fmt.Sprintf("%v = (%v)(%v)", arg.GoName(), arg.GoType(), arg.CName())
	}
}

func (arg *Argument) ConvertToC() string {
	if arg.Error {
		return fmt.Sprintf("var %v %v", arg.CName(), arg.CType())
	}

	ti := arg.TypeInfo()
	switch tag := ti.GetTag(); tag {
	case gi.TypeTagBoolean:
		return fmt.Sprintf("var %v C.gboolean\nif %v { %v = 1 }", arg.CName(), arg.GoName(), arg.CName())

	case gi.TypeTagUtf:
		var free string
		if arg.Info.GetOwnershipTransfer() == gi.TransferNothing {
			free = fmt.Sprintf("\ndefer C.free(unsafe.Pointer(%v))", arg.CName())
		}
		return fmt.Sprintf("%v := C.CString(%v)%v", arg.CName(), arg.GoName(), free)

	case gi.TypeTagInterface:
		return fmt.Sprintf("%v := (%v)(unsafe.Pointer(%v))", arg.CName(), arg.CType(), arg.GoName())

	default:
		return fmt.Sprintf("%v := (%v)(%v)", arg.CName(), arg.CType(), arg.GoName())
	}
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
