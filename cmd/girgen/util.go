package main

var (
	typeTagsGo = [...]string{
		"INVALID(Void)",      // Void
		"bool",               // Boolean
		"int8",               // Int8
		"uint8",              // Uint8
		"int16",              // Int16
		"uint16",             // Uint16
		"int32",              // Int32
		"uint32",             // Uint32
		"int64",              // Int64
		"uint64",             // Uint64
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
		"C.schar",            // Int8
		"C.uchar",            // Uint8
		"C.short",            // Int16
		"C.ushort",           // Uint16
		"C.int",              // Int32
		"C.uint",             // Uint32
		"C.longlong",         // Int64
		"C.size_t",           // Uint64
		"C.float",            // Float
		"C.double",           // Double
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
