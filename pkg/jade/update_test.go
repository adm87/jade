package jade_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adm87/jade/pkg/jade"
)

var errBoom = errors.New("boom")

func TestUpdateAll(t *testing.T) {
	var got []string
	record := func(name string, err error) jade.Updater {
		return jade.UpdateFunc(func(time.Duration) error {
			got = append(got, name)
			return err
		})
	}

	tests := []struct {
		name     string
		updaters []jade.Updater
		want     []string
		wantErr  error
	}{
		{"runs in order", []jade.Updater{record("a", nil), record("b", nil), record("c", nil)}, []string{"a", "b", "c"}, nil},
		{"skips nil", []jade.Updater{record("a", nil), nil, record("b", nil)}, []string{"a", "b"}, nil},
		{"stops at first error", []jade.Updater{record("a", nil), record("b", errBoom), record("c", nil)}, []string{"a", "b"}, errBoom},
		{"empty", nil, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got = nil
			err := jade.UpdateAll(tt.updaters...).Update(time.Millisecond)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFixedStep(t *testing.T) {
	const step = 16 * time.Millisecond
	ms := time.Millisecond

	tests := []struct {
		name      string
		dts       []time.Duration
		wantCalls int
	}{
		{"exact step", []time.Duration{16 * ms}, 1},
		{"carries leftover", []time.Duration{10 * ms, 10 * ms}, 1},
		{"catches up", []time.Duration{40 * ms}, 2},
		{"caps steps and drops the rest", []time.Duration{100 * ms, 10 * ms}, 5},
		{"ignores zero dt", []time.Duration{0}, 0},
		{"ignores negative dt", []time.Duration{-16 * ms, 16 * ms}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			u := jade.FixedStep(step, 5, jade.UpdateFunc(func(dt time.Duration) error {
				calls++
				assert.Equal(t, step, dt)
				return nil
			}))

			for _, dt := range tt.dts {
				require.NoError(t, u.Update(dt))
			}
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestFixedStepError(t *testing.T) {
	calls := 0
	u := jade.FixedStep(16*time.Millisecond, 5, jade.UpdateFunc(func(time.Duration) error {
		calls++
		return errBoom
	}))

	assert.ErrorIs(t, u.Update(32*time.Millisecond), errBoom)
	assert.Equal(t, 1, calls)
}

func TestFixedStepNil(t *testing.T) {
	assert.NoError(t, jade.FixedStep(16*time.Millisecond, 5, nil).Update(time.Second))
}
