package slotmap

import (
	"iter"
	"math"
	"slices"
)

// Every slot index has to fit in a uint32, and every slice length in an int.
const maxSlots = min(math.MaxUint32, math.MaxInt)

type slot struct {
	pos     uint32 // index into values while occupied
	next    uint32 // next free slot while free
	version uint32 // low bit is set while the slot holds a value
}

// Key finds a value in a Map. The zero Key is never valid.
type Key struct {
	index   uint32
	version uint32
}

// Pack fits k into a uint64.
func Pack(k Key) uint64 {
	return uint64(k.version)<<32 | uint64(k.index)
}

// Unpack turns a uint64 from Pack back into a Key.
func Unpack(packed uint64) Key {
	return Key{
		index:   uint32(packed & math.MaxUint32),
		version: uint32(packed >> 32),
	}
}

// Map holds values under keys that stop working once the value is deleted.
// Its zero value is ready to use.
type Map[T any] struct {
	slots   []slot
	values  []T    // packed so All only visits live values
	keys    []Key  // key of each value
	head    uint32 // first free slot, or 0 for none (slot 0 is never used)
	retired int    // slots whose version ran out
}

// New returns a Map with room for capacity values.
func New[T any](capacity int) *Map[T] {
	m := &Map[T]{}
	m.Grow(capacity)
	return m
}

// Len is the number of values in m.
func (m *Map[T]) Len() int {
	return len(m.values)
}

// Cap is the number of values m can hold before it grows.
func (m *Map[T]) Cap() int {
	return max(len(m.slots)-1-m.retired, 0)
}

// Insert adds value and returns its key. It panics if m would need more
// than math.MaxUint32 slots.
func (m *Map[T]) Insert(value T) Key {
	if m.head == 0 {
		m.Grow(growBy(len(m.slots)))
	}

	index := m.head
	sl := &m.slots[index]
	m.head = sl.next

	sl.version++
	sl.pos = uint32(len(m.values)) //nolint:gosec // bounded by maxSlots
	k := Key{index: index, version: sl.version}
	m.values = append(m.values, value)
	m.keys = append(m.keys, k)
	return k
}

// Get returns the value for k, or false if it was deleted.
func (m *Map[T]) Get(k Key) (T, bool) {
	sl, ok := m.lookup(k)
	if !ok {
		var zero T
		return zero, false
	}
	return m.values[sl.pos], true
}

// GetPtr returns a pointer to the value for k, or false if it was deleted.
// The pointer is only valid until the next Insert, Delete, Grow or Clear.
func (m *Map[T]) GetPtr(k Key) (*T, bool) {
	sl, ok := m.lookup(k)
	if !ok {
		return nil, false
	}
	return &m.values[sl.pos], true
}

// Set replaces the value for k. It returns false if k was deleted.
func (m *Map[T]) Set(k Key, value T) bool {
	sl, ok := m.lookup(k)
	if !ok {
		return false
	}
	m.values[sl.pos] = value
	return true
}

// Contains reports whether k still has a value.
func (m *Map[T]) Contains(k Key) bool {
	_, ok := m.lookup(k)
	return ok
}

// Delete removes the value for k and returns it, or false if it was already deleted.
func (m *Map[T]) Delete(k Key) (T, bool) {
	var zero T
	sl, ok := m.lookup(k)
	if !ok {
		return zero, false
	}
	old := m.values[sl.pos]

	// Move the last value into the hole so values stay packed.
	last := len(m.values) - 1
	m.values[sl.pos] = m.values[last]
	m.keys[sl.pos] = m.keys[last]
	m.slots[m.keys[sl.pos].index].pos = sl.pos
	m.values[last] = zero
	m.values = m.values[:last]
	m.keys = m.keys[:last]

	sl.version++
	if sl.version == math.MaxUint32-1 {
		// One more use would wrap the version and bring old keys back.
		m.retired++
		return old, true
	}
	sl.next = m.head
	m.head = k.index
	return old, true
}

// Clear deletes every value. Every existing key stops working.
func (m *Map[T]) Clear() {
	for i := range slices.Backward(m.values) {
		m.Delete(m.keys[i])
	}
}

// All yields every key and value in no particular order.
// Deleting the current key during All is safe.
func (m *Map[T]) All() iter.Seq2[Key, T] {
	return func(yield func(Key, T) bool) {
		// Walk backwards, so deleting the current value only moves one already
		// visited. Read m.values fresh each time, since an Insert may have grown it.
		for i := range slices.Backward(m.values) {
			if i >= len(m.values) {
				continue
			}
			if !yield(m.keys[i], m.values[i]) {
				return
			}
		}
	}
}

// Values yields every value in no particular order.
func (m *Map[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		// Walk backwards for the same reasons as All.
		for i := range slices.Backward(m.values) {
			if i >= len(m.values) {
				continue
			}
			if !yield(m.values[i]) {
				return
			}
		}
	}
}

// Grow adds room for n more values. It panics if m would need more than
// math.MaxUint32 slots.
func (m *Map[T]) Grow(n int) {
	if n <= 0 {
		return
	}
	if len(m.slots) == 0 {
		m.slots = append(m.slots, slot{})
	}
	if n > maxSlots-len(m.slots) {
		panic("slotmap: too many slots")
	}

	oldLen := uint32(len(m.slots)) //nolint:gosec // bounded by maxSlots
	m.slots = append(m.slots, make([]slot, n)...)
	newTail := uint32(len(m.slots) - 1) //nolint:gosec // bounded by maxSlots

	for i := oldLen; i < newTail; i++ {
		m.slots[i].next = i + 1
	}

	// Put the new free slots in front of the old ones.
	m.slots[newTail].next = m.head
	m.head = oldLen

	m.values = slices.Grow(m.values, m.Cap()-len(m.values))
	m.keys = slices.Grow(m.keys, m.Cap()-len(m.keys))
}

// lookup returns the slot for k, or false if k is invalid or its value was deleted.
func (m *Map[T]) lookup(k Key) (*slot, bool) {
	if k.index == 0 || uint64(k.index) >= uint64(len(m.slots)) {
		return nil, false
	}
	sl := &m.slots[k.index]
	if sl.version != k.version || sl.version&1 == 0 {
		return nil, false
	}
	return sl, true
}

// growBy doubles the slots, but never past maxSlots.
func growBy(slots int) int {
	return max(min(slots, maxSlots-slots), 1)
}
