package util

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"slices"
	"strings"
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

func JoinSeq[T any](seq iter.Seq[T], sep string) string {
	return strings.Join(slices.Collect(func(yield func(string) bool) {
		for v := range seq {
			yield(fmt.Sprint(v))
		}
	}), sep)
}

func JoinPairs[T1, T2 any](seq iter.Seq2[T1, T2], psep, sep string) string {
	return strings.Join(slices.Collect(func(yield func(string) bool) {
		for v1, v2 := range seq {
			yield(fmt.Sprintf("%v%v%v", v1, psep, v2))
		}
	}), sep)
}
