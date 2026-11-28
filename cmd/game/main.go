package main

import (
	"fmt"
	"os"

	"github.com/adm87/jade/cmd/game/internal/game"
)

func main() {
	if err := game.Boot(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
