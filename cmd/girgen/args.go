package main

import (
	"iter"

	"deedles.dev/gir/gi"
	"deedles.dev/gir/internal/util"
)

type Arguments struct {
	Gen  *Generator
	Args []*Argument
}

func (args *Arguments) Callable() *gi.CallableInfo {
	return args.Gen.Element.(interface{ AsGICallableInfo() *gi.CallableInfo }).AsGICallableInfo()
}

func (args *Arguments) GoInput() string {
	return util.JoinPairs(args.goInput(), " ", ", ")
}

func (args *Arguments) goInput() iter.Seq2[string, string] {
}

func (args *Arguments) CInput() string {
	return util.JoinSeq(args.cInput(), ", ")
}

func (args *Arguments) cInput() iter.Seq[string] {
}

func (args *Arguments) ConvertToC() string {
}

type Argument struct {
	Index uint
	Info  *gi.ArgInfo
}
