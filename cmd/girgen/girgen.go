package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"deedles.dev/gir/gi"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: girgen < config > output.go")
	}
	flag.Parse()

	config, err := ParseConfig(os.Stdin)
	if err != nil {
		slog.Error("failed to parse config", "err", err)
		os.Exit(1)
	}

	r := gi.RepositoryNew()

	tl, err := r.Require(config.Namespace, config.Version, gi.RepositoryLoadFlagNone)
	if err != nil {
		slog.Error("failed to open typelib", "namespace", config.Namespace, "version", config.Version, "err", err)
		os.Exit(1)
	}
	defer tl.Unref()

	for info := range r.GetInfos(config.Namespace) {
		fmt.Printf("%q\n", info.GetName())
		info.Unref()
	}
}
