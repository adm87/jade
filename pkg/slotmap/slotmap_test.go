package slotmap_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adm87/jade/pkg/slotmap"
)

// index is the slot index, the low half of a packed key.
func index(k slotmap.Key) uint32 {
	return uint32(slotmap.Pack(k) & math.MaxUint32)
}

func TestInsertGet(t *testing.T) {
	m := slotmap.New[string](2)
	a := m.Insert("a")
	b := m.Insert("b")

	got, ok := m.Get(a)
	assert.True(t, ok)
	assert.Equal(t, "a", got)

	got, ok = m.Get(b)
	assert.True(t, ok)
	assert.Equal(t, "b", got)

	assert.Equal(t, 2, m.Len())
}

func TestZeroValue(t *testing.T) {
	var m slotmap.Map[int]
	assert.Equal(t, 0, m.Len())
	assert.Equal(t, 0, m.Cap())
	assert.False(t, m.Contains(slotmap.Key{}))

	k := m.Insert(1)
	got, ok := m.Get(k)
	assert.True(t, ok)
	assert.Equal(t, 1, got)
	assert.NotZero(t, index(k), "slot 0 is reserved")
}

func TestInvalidKey(t *testing.T) {
	m := slotmap.New[int](1)
	m.Insert(1)

	tests := []struct {
		name string
		key  slotmap.Key
	}{
		{"zero", slotmap.Key{}},
		{"out of range", slotmap.Unpack(1<<32 | 5)},
		{"wrong version", slotmap.Unpack(3<<32 | 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := m.Get(tt.key)
			assert.False(t, ok, "Get")
			_, ok = m.GetPtr(tt.key)
			assert.False(t, ok, "GetPtr")
			assert.False(t, m.Set(tt.key, 2), "Set")
			assert.False(t, m.Contains(tt.key), "Contains")
			_, ok = m.Delete(tt.key)
			assert.False(t, ok, "Delete")
		})
	}
}

func TestDeleteStalesKey(t *testing.T) {
	m := slotmap.New[int](1)
	k := m.Insert(7)

	got, ok := m.Delete(k)
	require.True(t, ok)
	assert.Equal(t, 7, got)

	assert.False(t, m.Contains(k))
	_, ok = m.Delete(k)
	assert.False(t, ok, "second Delete")

	k2 := m.Insert(8)
	assert.Equal(t, index(k), index(k2), "reuses the slot")
	assert.NotEqual(t, k, k2, "with a new version")
	assert.False(t, m.Contains(k), "old key stays stale")
	assert.Equal(t, 1, m.Len())
}

func TestDeleteKeepsOtherKeys(t *testing.T) {
	m := slotmap.New[string](3)
	a := m.Insert("a")
	b := m.Insert("b")
	c := m.Insert("c")

	m.Delete(a)

	got, ok := m.Get(b)
	assert.True(t, ok)
	assert.Equal(t, "b", got)

	got, ok = m.Get(c)
	assert.True(t, ok)
	assert.Equal(t, "c", got)
}

func TestSetAndGetPtr(t *testing.T) {
	var m slotmap.Map[int]
	k := m.Insert(1)

	require.True(t, m.Set(k, 2))
	p, ok := m.GetPtr(k)
	require.True(t, ok)
	assert.Equal(t, 2, *p)

	*p = 3
	got, _ := m.Get(k)
	assert.Equal(t, 3, got)
}

func TestGrowKeepsKeys(t *testing.T) {
	m := slotmap.New[int](0)
	require.Equal(t, 0, m.Cap())

	keys := make([]slotmap.Key, 100)
	for i := range keys {
		keys[i] = m.Insert(i)
	}
	for i, k := range keys {
		got, ok := m.Get(k)
		assert.True(t, ok)
		assert.Equal(t, i, got)
	}
	assert.GreaterOrEqual(t, m.Cap(), 100)
}

func TestGrowWithFreeSlots(t *testing.T) {
	m := slotmap.New[int](2)
	m.Delete(m.Insert(1))
	m.Grow(2)

	seen := map[uint32]bool{}
	for range 4 {
		seen[index(m.Insert(0))] = true
	}
	assert.Len(t, seen, 4, "distinct slots")
	assert.Equal(t, 4, m.Cap())
}

func TestGrowPastMaxPanics(t *testing.T) {
	m := slotmap.New[struct{}](0)
	assert.Panics(t, func() { m.Grow(int(^uint(0) >> 1)) })
}

func TestClear(t *testing.T) {
	m := slotmap.New[int](3)
	keys := []slotmap.Key{m.Insert(1), m.Insert(2), m.Insert(3)}

	m.Clear()
	assert.Equal(t, 0, m.Len())
	for _, k := range keys {
		assert.False(t, m.Contains(k))
	}

	k := m.Insert(4)
	assert.NotContains(t, keys, k, "new keys don't match old ones")
	assert.Equal(t, 3, m.Cap(), "slots are reused")
}

func TestAll(t *testing.T) {
	m := slotmap.New[int](4)
	for i := range 4 {
		m.Insert(i)
	}
	m.Delete(m.Insert(99))

	sum := 0
	for k, v := range m.All() {
		got, _ := m.Get(k)
		assert.Equal(t, got, v, "All matches Get")
		sum += v
	}
	assert.Equal(t, 0+1+2+3, sum)

	calls := 0
	for range m.All() {
		calls++
		break
	}
	assert.Equal(t, 1, calls, "stops on break")
}

func TestDeleteDuringAll(t *testing.T) {
	m := slotmap.New[int](5)
	for i := range 5 {
		m.Insert(i)
	}

	seen := map[int]bool{}
	for k, v := range m.All() {
		assert.False(t, seen[v], "yielded %d twice", v)
		seen[v] = true
		m.Delete(k)
	}
	assert.Len(t, seen, 5)
	assert.Equal(t, 0, m.Len())
}

func TestAllSeesSetAfterGrowth(t *testing.T) {
	m := slotmap.New[int](3)
	a := m.Insert(1)
	m.Insert(2)
	m.Insert(3)

	grown := false
	var got int
	for k, v := range m.All() {
		if !grown {
			m.Insert(4) // full, so this reallocates
			m.Set(a, 100)
			grown = true
		}
		if k == a {
			got = v
		}
	}
	assert.Equal(t, 100, got)
}

func TestPackRoundTrip(t *testing.T) {
	m := slotmap.New[int](1)
	k := m.Insert(1)
	assert.Equal(t, k, slotmap.Unpack(slotmap.Pack(k)))
}

func TestValues(t *testing.T) {
	m := slotmap.New[int](4)
	for i := range 4 {
		m.Insert(i)
	}
	m.Delete(m.Insert(99))

	var got []int
	for v := range m.Values() {
		got = append(got, v)
	}
	assert.ElementsMatch(t, []int{0, 1, 2, 3}, got)

	calls := 0
	for range m.Values() {
		calls++
		break
	}
	assert.Equal(t, 1, calls, "stops on break")
}
