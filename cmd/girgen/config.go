package main

import (
	"fmt"
	"io"
	"strings"

	"deedles.dev/gir/internal/util"
)

type Config struct {
	Package    string
	Namespace  string
	Version    string
	Includes   []string
	PkgConfigs []string
}

func ParseConfig(r io.Reader) (*Config, error) {
	var config Config
	for line, err := range util.Lines(r) {
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		err := config.run(strings.Fields(line))
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func (config *Config) run(directive []string) error {
	switch directive[0] {
	case "package":
		return assertArgs(directive, 1, 1, func() {
			config.Package = directive[1]
		})

	case "namespace":
		return assertArgs(directive, 1, 1, func() {
			config.Namespace = directive[1]
		})

	case "version":
		return assertArgs(directive, 1, 1, func() {
			config.Version = directive[1]
		})

	case "include":
		return assertArgs(directive, 1, -1, func() {
			config.Includes = append(config.Includes, directive[1:]...)
		})

	case "pkg-config":
		return assertArgs(directive, 1, -1, func() {
			config.PkgConfigs = append(config.PkgConfigs, directive[1:]...)
		})

	default:
		return fmt.Errorf("unknown directive %q", directive[0])
	}
}

func assertArgs(directive []string, min, max int, f func()) error {
	args := len(directive[1:])
	switch {
	case min == max && args != min:
		return fmt.Errorf("directive %q expects %v argument(s) but was given %v", directive[0], min, args)
	case args < min:
		return fmt.Errorf("directive %q expects at least %v argument(s) but was given %v", directive[0], min, args)
	case max > 0 && args > max:
		return fmt.Errorf("directive %q expects at most %v argument(s) but was given %v", directive[0], max, args)
	}

	f()
	return nil
}
