package ttt

import "testing"

func Test_GameSetup(t *testing.T) {

}


func test() {
	m := NewGame()
	x := Mover(m.GameMap, FillX)
	o := Mover(m.GameMap, FillO)

	o(1, 1)
	x(0, 2)
	o(2, 2)
	x(0, 1)
	o(3, 3)
	o(-1, 2)
	o(-2, 2)
	x(-2, 1)
	x(-3, -1)
	o(-5, -5)
	x(5, 4)
	o(7, 4)
	m.draw()

}