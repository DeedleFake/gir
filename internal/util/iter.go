package util

import "iter"

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
