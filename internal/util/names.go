package util

import (
	"fmt"
	"strings"
	"unicode"
)

func ToCamelCase(name string) string {
	var buf strings.Builder
	for part := range strings.SplitSeq(name, "_") {
		buf.WriteString(strings.ToUpper(part[:1]))
		buf.WriteString(strings.ToLower(part[1:]))
	}
	return buf.String()
}

func ToSnakeCase(name string) string {
	var buf strings.Builder
	for i, c := range name {
		if unicode.IsUpper(c) && i != 0 {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(c))
	}
	return buf.String()
}

func ParseCPrefix(cprefix string) string {
	prefix, _, _ := strings.Cut(cprefix, ",")
	return prefix
}

func MethodName(prefix, tname, mname string) string {
	return fmt.Sprintf("%v_%v_%v", strings.ToLower(prefix), ToSnakeCase(tname), mname)
}
