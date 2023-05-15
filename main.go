package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"guitarboard/game"
	"guitarboard/img"
	"log"
)

func main() {
	img.InitImg()
	ebiten.SetWindowSize(1200, 400)
	ebiten.SetWindowTitle("BassBoard")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
