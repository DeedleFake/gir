package main

var (
	typeTagsGo = [...]string{
		"INVALID(Void)",      // Void
		"bool",               // Boolean
		"int8",               // Int
		"uint8",              // Uint
		"int16",              // Int
		"uint16",             // Uint
		"int32",              // Int
		"uint32",             // Uint
		"int64",              // Int
		"uint64",             // Uint
		"float32",            // Float
		"float64",            // Double
		"INVALID(GType)",     // GType
		"string",             // Utf
		"string",             // Filename
		"INVALID(Array)",     // Array
		"INVALID(Interface)", // Interface
		"INVALID(GList)",     // GList
		"INVALID(GSList)",    // GSList
		"INVALID(GHash)",     // GHash
		"INVALID(Error)",     // Error
		"INVALID(Unichar)",   // Unichar
	}

	typeTagsC = [...]string{
		"INVALID(Void)",      // Void
		"C.gboolean",         // Boolean
		"C.int8",             // Int
		"C.uint8",            // Uint
		"C.int16",            // Int
		"C.uint16",           // Uint
		"C.int32",            // Int
		"C.uint32",           // Uint
		"C.int64",            // Int
		"C.uint64",           // Uint
		"C.float32",          // Float
		"C.float64",          // Double
		"C.GType",            // Gtype
		"*C.char",            // Utf
		"*C.char",            // Filename
		"INVALID(Array)",     // Array
		"INVALID(Interface)", // Interface
		"INVALID(GList)",     // GList
		"INVALID(GSList)",    // GSList
		"INVALID(GHash)",     // GHash
		"INVALID(Error)",     // Error
		"INVALID(Unichar)",   // Unichar
	}
)

type parentInfo interface {
	GetNamespace() string
	GetName() string
}

type typeInstanceParentInfo struct{}

func (typeInstanceParentInfo) GetNamespace() string { return "GObject" }
func (typeInstanceParentInfo) GetName() string      { return "TypeInstance" }
