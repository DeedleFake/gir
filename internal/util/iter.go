package util

import (
	"bufio"
	"io"
	"iter"
)

func Pairs[T any](seq iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		var prev *T
		for v := range seq {
			if prev == nil {
				prev = &v
				continue
			}

			if !yield(*prev, v) {
				return
			}
			prev = nil
		}
		if prev != nil {
			panic("odd length of paris")
		}
	}
}

func Lines(r io.Reader) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if !yield(s.Text(), nil) {
				return
			}
		}
		if s.Err() != nil {
			yield("", s.Err())
		}
	}
}
