package slotmap_test

import (
	"testing"

	"github.com/adm87/jade/pkg/slotmap"
)

const benchSize = 10_000

// filled returns a Map holding benchSize values and their keys.
func filled() (*slotmap.Map[int], []slotmap.Key) {
	m := slotmap.New[int](benchSize)
	keys := make([]slotmap.Key, benchSize)
	for i := range keys {
		keys[i] = m.Insert(i)
	}
	return m, keys
}

func BenchmarkInsert(b *testing.B) {
	b.Run("growing", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var m slotmap.Map[int]
			for i := range benchSize {
				m.Insert(i)
			}
		}
	})
	b.Run("preallocated", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			m := slotmap.New[int](benchSize)
			for i := range benchSize {
				m.Insert(i)
			}
		}
	})
}

func BenchmarkGet(b *testing.B) {
	m, keys := filled()
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		m.Get(keys[i%benchSize])
		i++
	}
}

func BenchmarkGetPtr(b *testing.B) {
	m, keys := filled()
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		p, _ := m.GetPtr(keys[i%benchSize])
		*p++
		i++
	}
}

// BenchmarkChurn deletes a value and inserts a new one, which reuses the freed slot.
func BenchmarkChurn(b *testing.B) {
	m, keys := filled()
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		j := i % benchSize
		m.Delete(keys[j])
		keys[j] = m.Insert(i)
		i++
	}
}

func BenchmarkAll(b *testing.B) {
	m, _ := filled()
	b.ReportAllocs()
	for b.Loop() {
		sum := 0
		for _, v := range m.All() {
			sum += v
		}
		_ = sum
	}
}

func BenchmarkValues(b *testing.B) {
	m, _ := filled()
	b.ReportAllocs()
	for b.Loop() {
		sum := 0
		for v := range m.Values() {
			sum += v
		}
		_ = sum
	}
}
