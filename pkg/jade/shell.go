package jade

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// maxDelta caps dt so a long pause doesn't produce one huge step.
const maxDelta = time.Second

type shell struct {
	last     time.Time
	drawer   Drawer
	updater  Updater
	layouter Layouter
}

func (s *shell) Update() error {
	return s.updater.Update(s.dt())
}

func (s *shell) Draw(screen *ebiten.Image) {
	s.drawer.Draw(screen)
}

// Layout is required by ebiten.Game, but Ebitengine calls LayoutF instead when it exists.
func (s *shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	w, h := s.LayoutF(float64(outsideWidth), float64(outsideHeight))
	return int(w), int(h)
}

func (s *shell) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
	return s.layouter.Layout(outsideWidth, outsideHeight)
}

// dt returns the time since the previous call, capped at maxDelta. It returns 0 on the first call.
func (s *shell) dt() time.Duration {
	now := time.Now()
	last := s.last
	s.last = now

	if last.IsZero() {
		return 0
	}

	return min(now.Sub(last), maxDelta)
}
