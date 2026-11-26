//go:build integration

package jade_test

import (
	"os"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adm87/jade/pkg/jade"
)

// Results of the one Run per process, checked by the tests below.
var (
	runErr error
	dts    []time.Duration
	draws  int
)

// TestMain calls Run because Ebitengine must start on the main thread,
// and only TestMain runs there.
func TestMain(m *testing.M) {
	runErr = jade.Run(
		jade.WithUpdater(jade.UpdateFunc(func(dt time.Duration) error {
			dts = append(dts, dt)
			if draws > 0 && len(dts) > 2 || len(dts) > 600 {
				return ebiten.Termination
			}
			return nil
		})),
		jade.WithDrawer(jade.DrawFunc(func(*ebiten.Image) { draws++ })),
	)
	os.Exit(m.Run())
}

func TestIntegrationRun(t *testing.T) {
	require.NoError(t, runErr, "Termination should end Run cleanly")
	require.Greater(t, len(dts), 2)
	assert.Positive(t, draws)
	assert.Zero(t, dts[0], "first tick")
	assert.Positive(t, dts[1], "measured after the first tick")
}
