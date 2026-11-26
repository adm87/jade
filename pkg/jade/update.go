package jade

import (
	"time"
)

// Updater advances game state by one tick.
type Updater interface {
	Update(dt time.Duration) error
}

// UpdateFunc lets an ordinary function be used as an Updater.
type UpdateFunc func(dt time.Duration) error

// Update calls f(dt).
func (f UpdateFunc) Update(dt time.Duration) error {
	return f(dt)
}

// UpdatePass wraps an Updater to control whether, when, or how it runs.
type UpdatePass func(Updater) Updater

// UpdateAll runs updaters in order, stopping at the first error.
// Nil updaters are skipped.
func UpdateAll(updaters ...Updater) Updater {
	return UpdateFunc(func(dt time.Duration) error {
		for _, u := range updaters {
			if u == nil {
				continue
			}
			if err := u.Update(dt); err != nil {
				return err
			}
		}
		return nil
	})
}

// FixedStep runs updater in steps of exactly step, catching up on accumulated time.
// At most maxSteps run per tick, and any time beyond that is dropped.
func FixedStep(step time.Duration, maxSteps int, updater Updater) Updater {
	if updater == nil {
		return UpdateFunc(func(time.Duration) error { return nil })
	}

	var acc time.Duration
	return UpdateFunc(func(dt time.Duration) error {
		if dt <= 0 {
			return nil
		}

		acc += dt
		for n := 0; acc >= step; n++ {
			if n >= maxSteps {
				acc = 0
				break
			}
			if err := updater.Update(step); err != nil {
				return err
			}
			acc -= step
		}

		return nil
	})
}
