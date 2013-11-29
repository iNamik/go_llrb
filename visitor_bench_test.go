package llrb

import (
	"math/rand"
	"testing"
)

import (
	"github.com/iNamik/go_bst/visitor"
)

func Benchmark_Visit1_Small(b *testing.B) {
	const SIZE = 10
	const ITERATIONS = 10000
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		r := randomTree(SIZE)
		b.StartTimer()
		for i := 0; i < ITERATIONS; i++ {
			n := rand.Intn(SIZE)
			r.Visit(n, func(v_ interface{}, found bool) (interface{}, visitor.Action) {
				if found {
					return nil, visitor.REMOVE
				} else {
					return n, visitor.INSERT
				}
			})
		}
	}
}

func Benchmark_Visit1_Large(b *testing.B) {
	const SIZE = 100
	const ITERATIONS = 1000000
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		r := randomTree(SIZE)
		b.StartTimer()
		for i := 0; i < ITERATIONS; i++ {
			n := rand.Intn(SIZE)
			r.Visit(n, func(v_ interface{}, found bool) (interface{}, visitor.Action) {
				if found {
					return nil, visitor.REMOVE
				} else {
					return n, visitor.INSERT
				}
			})
		}
	}
}

func Benchmark_Visit2_Small(b *testing.B) {
	const SIZE = 10
	const ITERATIONS = 10000
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		r := randomTree(SIZE)
		b.StartTimer()
		for i := 0; i < ITERATIONS; i++ {
			n := rand.Intn(SIZE)
			if _, found := r.Get(n); found {
				r.Remove(n)
			} else {
				r.ReplaceOrInsert(n, n)
			}
		}
	}
}

func Benchmark_Visit2_Large(b *testing.B) {
	const SIZE = 100
	const ITERATIONS = 1000000
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		r := randomTree(SIZE)
		b.StartTimer()
		for i := 0; i < ITERATIONS; i++ {
			n := rand.Intn(SIZE)
			if _, found := r.Get(n); found {
				r.Remove(n)
			} else {
				r.ReplaceOrInsert(n, n)
			}
		}
	}
}
