package main

import jade "github.com/adm87/jade/cmd/jade-game/internal"

func main() {
	if err := jade.Boot(); err != nil {
		println("error:", err.Error())
	}
}
