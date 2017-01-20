package ttt

import "fmt"


type Fill int

const (
	FillNil Fill = 0
	FillX   Fill = 1
	FillO   Fill = 2
)
type Status int

const (
	Waiting = 0
	Ready = 1
	Finished = 2
)

type (
	Point struct{ X, Y int }
	Node  struct {
		Point
		Content Fill
	}

	GameMap map[string]*Node
	Game struct{
		GameMap
		Status Status
		Winner Fill
		Next Fill
	}
)

func NewGame() *Game {
	g := new(Game)
	m := make(GameMap)
	m[Point{0, 0}.String()] = NewNode(0, 0, FillX)
	g.GameMap = m
	return g
}


func (game *Game) GameLoop(c chan Point, n chan struct{}) {

	x := Mover(game.GameMap, FillX)
	o := Mover(game.GameMap, FillO)

	game.Next= FillO
	for p := range c {

		if game.Next == FillO {
			o(p.X, p.Y)
			game.Next = FillX
		} else {
			x(p.X, p.Y)
			game.Next = FillO
		}

		if f := game.TestVictory(); f != FillNil {
			fmt.Println("Victory for" , f)
			game.Status = Finished
			game.Winner = f
		}
		n <- struct{}{}
	}

}



func NewNode(x, y int, fill Fill) *Node {
	return &Node{Point{x, y}, fill}
}

func Mover(m GameMap, fill Fill) func(int, int) {
	return func(x, y int) {
		m.Move(x, y, fill)
	}
}

func (m GameMap) Move(x, y int, fill Fill) {
	m[Point{x, y}.String()] = NewNode(x, y, fill)
}


func(m GameMap) TestVictory () Fill {
	for _, n := range m {
		for _, d := range m.directions(*n) {
			fmt.Println("tesng ", n, d)
			if f := m.testDirection(d, 1); f != FillNil{
				return f
			}
		}
	}
	return FillNil

}

func(m GameMap) testDirection(d Direction, s int) Fill {
	fmt.Println("testing direction ", d, s)
	n := m[d.Point.String()]
	dirs := directions(n.Point)
	td := dirs[d.i]
	if v, ok := m[ td.Point.String() ]; ok && v.Content == n.Content {
		if s == 3 {
			return n.Content
		}
		return m.testDirection(td, s+1)
	}
	return FillNil
}

type Direction struct {
	Point
	i int
}
func (m GameMap) directions(n Node) []Direction {
	res := make([]Direction, 0)

	for _, d := range directions(n.Point) {
		if v, ok := m[ d.Point.String() ]; ok && n.Content == v.Content {
			res = append(res, d)
		}
	}
	return res
}

func directions(p Point) []Direction {
	res := make([]Direction, 8)
	res[0] = Direction{Point{p.X, p.Y-1}, 0}
	res[1] = Direction{Point{p.X+1, p.Y-1}, 1}
	res[2] = Direction{Point{p.X+1, p.Y}, 2}
	res[3] = Direction{Point{p.X+1, p.Y+1}, 3}
	res[4] = Direction{Point{p.X, p.Y+1}, 4}
	res[5] = Direction{Point{p.X-1, p.Y+1}, 5}
	res[6] = Direction{Point{p.X-1, p.Y}, 6}
	res[7] = Direction{Point{p.X-1, p.Y-1}, 7}
	return res
}

type Nodes []Node

func (n Nodes) Len() int      { return len(n) }
func (n Nodes) Swap(a, b int) { n[a], n[b] = n[b], n[a] }
func (n Nodes) Less(a, b int) bool {
	if n[a].X < n[b].Y {
		return true
	}
	if n[a].X > n[b].Y {
		return false
	}
	if n[a].Y < n[b].Y {
		return true
	}
	return false
}

func (m GameMap) draw() {
	fmt.Println(m)
	var minX, maxX, minY, maxY int
	for _, p := range m {

		if minX > p.X {
			minX = p.X
		}
		if minY > p.Y {
			minY = p.Y
		}
		if maxX < p.X {
			maxX = p.X
		}
		if maxY < p.Y {
			maxY = p.Y
		}

	}
	for y := minY - 2; y <= maxY+2; y++ {

		for x := minX - 2; x <= maxX+2; x++ {
			n, ok := m[Point{x, y}.String()]
			if ok {
				fmt.Print(n)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")

	}
}

func (n Node) String() string {
	if n.Content == FillO {
		return "O"
	}
	return "X"
}

func (p Point) String() string {
	return fmt.Sprintf("%d.%d", p.X, p.Y)
}
func(d Direction) String() string{

	return fmt.Sprintf("%d.%d = %d", d.X, d.Y, d.i)
}