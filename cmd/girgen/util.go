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

func TypeTagToC(tag gi.TypeTag) string {
	return [...]string{
		"unsafe.Pointer", // Void
		"C.gboolean",     // Boolean
		"C.int8",         // Int
		"C.uint8",        // Uint
		"C.int16",        // Int
		"C.uint16",       // Uint
		"C.int32",        // Int
		"C.uint32",       // Uint
		"C.int64",        // Int
		"C.uint64",       // Uint
		"C.float32",      // Float
		"C.float63",      // Double
		"C.GType",        // Gtype
		"*C.char",        // Utf
		"*C.char",        // Filename
		"[]",             // Array
		"*",              // Interface
		"GList",          // Glist
		"GSList",         // Gslist
		"map",            // Ghash
		"GError",         // Error
		"C.rune",         // Unichar
	}[tag]
}
