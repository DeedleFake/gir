package main

import (
	"fmt"

	"deedles.dev/gir/gi"
)

func main() {
	r := gi.RepositoryNew()
	tl, err := r.Require("GIRepository", "3.0", gi.RepositoryLoadFlagNone)
	if err != nil {
		panic(err)
	}
	defer tl.Unref()

	for info := range r.GetInfos("GIRepository") {
		defer info.Unref()
		fmt.Println(info.GetName())
		if info, ok := gi.TypeObjectInfo.Check(info.AsGTypeInstance()); ok {
			for method := range info.GetMethods() {
				defer method.Unref()
				fmt.Printf("\t%v\n", method.GetName())
			}
		}
	}
}
