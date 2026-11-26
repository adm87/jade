package main

import game "github.com/adm87/jade/cmd/game/internal"

func main() {
	if err := game.Boot(); err != nil {
		println("error:", err.Error())
	}
}
