package main

import (
	"fmt"
	"io"
	"strings"

	"deedles.dev/gir/internal/util"
)

type Config struct {
	Package   string
	Namespace string
	Version   string
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
		return assertArgs(directive, 1, func() {
			config.Package = directive[1]
		})

	case "namespace":
		return assertArgs(directive, 1, func() {
			config.Namespace = directive[1]
		})

	case "version":
		return assertArgs(directive, 1, func() {
			config.Version = directive[1]
		})

	default:
		return fmt.Errorf("unknown directive %q", directive[0])
	}
}

func assertArgs(directive []string, args int, f func()) error {
	if len(directive[1:]) != args {
		return fmt.Errorf("directive %q expects %v argument(s) but was given %v", directive[0], args, len(directive[1:]))
	}

	f()
	return nil
}
