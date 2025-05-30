package util

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) AddAll(vals iter.Seq[T]) {
	maps.Insert(s, func(yield func(T, struct{}) bool) {
		vals(func(v T) bool { return yield(v, struct{}{}) })
	})
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}
