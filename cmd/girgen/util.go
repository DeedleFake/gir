package main

import (
	"deedles.dev/gir/gi"
)

func TypeTagToGo(tag gi.TypeTag) string {
	return [...]string{
		"unsafe.Pointer", // Void
		"bool",           // Boolean
		"int8",           // Int
		"uint8",          // Uint
		"int16",          // Int
		"uint16",         // Uint
		"int32",          // Int
		"uint32",         // Uint
		"int64",          // Int
		"uint64",         // Uint
		"float32",        // Float
		"float63",        // Double
		"g.Type",         // Gtype
		"string",         // Utf
		"string",         // Filename
		"[]",             // Array
		"*",              // Interface
		"g.List",         // Glist
		"g.Slist",        // Gslist
		"map",            // Ghash
		"g.Error",        // Error
		"rune",           // Unichar
	}[tag]
}
