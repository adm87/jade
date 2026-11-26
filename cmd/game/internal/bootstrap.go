package internal

import "github.com/adm87/jade/pkg/jade"

// Boot initializes the sandbox game.
func Boot() error {
	return jade.Run()
}
