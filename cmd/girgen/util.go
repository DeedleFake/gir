package main

import "deedles.dev/gir/gi"

func TypeTagToGo(tag gi.TypeTag) string {
	return [...]string{
		"unsafe.Pointer",
		"bool",
		"int8",
		"uint8",
		"int16",
		"uint16",
		"int32",
		"uint32",
		"int64",
		"uint64",
		"float32",
		"float64",
		"g.Type",
		"string",
		"string",
		"[]",
		"*",
		"g.List",
		"g.SList",
		"map",
		"g.Error",
		"rune",
	}[tag]
}
