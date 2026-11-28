package slotmap

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrowByStopsAtMaxSlots(t *testing.T) {
	tests := []struct {
		slots, want int
	}{
		{0, 1},
		{4, 4},
		{maxSlots - 3, 3},
		{maxSlots, 1}, // Grow(1) then panics
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, growBy(tt.slots), "growBy(%d)", tt.slots)
	}
}

func TestDeleteRetiresWornOutSlot(t *testing.T) {
	m := New[int](1)
	k := m.Insert(1)

	k.version = math.MaxUint32 - 2
	m.slots[k.index].version = k.version
	_, ok := m.Delete(k)
	require.True(t, ok)

	assert.Equal(t, 0, m.Cap(), "the only slot is retired")
	assert.NotEqual(t, k.index, m.Insert(2).index, "a retired slot is never reused")
}
