package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Cell struct {
	Alive    bool
	Neigbors int
}

type Board struct {
	Width  int
	Height int
	Cells  [][]Cell
}

func NewBoard(width, height int) *Board {
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
	}
	return &Board{width, height, cells}
}

func (b *Board) Update() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.Cells[y][x].Neigbors = b.countCellNeighbors(x, y)
		}
	}

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if b.Cells[y][x].Neigbors < 2 || b.Cells[y][x].Neigbors > 3 {
				b.Cells[y][x].Alive = false
			} else if b.Cells[y][x].Neigbors == 3 {
				b.Cells[y][x].Alive = true
			}

		}
	}
}

func (b *Board) debugPrint() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			symbol := "[ ]"
			if b.Cells[y][x].Alive {
				symbol = "[#]"
			}
			fmt.Print(symbol, " ")
		}
		fmt.Print("\n")
	}
}

func (b *Board) countCellNeighbors(curX, curY int) (aliveNeighbors int) {
	neighbors := 0
	for y := -1; y <= 1; y++ {
		if curY+y < 0 || curY+y >= b.Height {
			continue
		}
		for x := -1; x <= 1; x++ {
			if x == 0 && y == 0 {
				continue
			}
			if curX+x < 0 || curX+x >= b.Width {
				continue
			}
			if b.Cells[curY+y][curX+x].Alive {
				neighbors++
			}
		}
	}
	return neighbors
}

type Game struct {
	width           int
	height          int
	lastFrameTime   time.Time
	accumulatedTime time.Duration
	gameSpeed       time.Duration

	zoom         float64
	offsetX      int
	offsetY      int
	lastMouseX   int
	lastMouseY   int
	mousePressed bool

	paused bool
	board  Board
}

func NewGame(width, height int) *Game {
	return &Game{
		width:         width,
		height:        height,
		zoom:          1.0,
		board:         *NewBoard(width, height),
		lastFrameTime: time.Now(),
		gameSpeed:     100 * time.Millisecond,
		paused:        true,
	}
}

func (g *Game) Update() error {
	deltaTime := time.Now().Sub(g.lastFrameTime)
	g.lastFrameTime = time.Now()
	g.accumulatedTime += deltaTime

	_, dy := ebiten.Wheel()
	g.zoom = max(1, g.zoom+dy)

	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if g.mousePressed {
			dx := x - g.lastMouseX
			dy := y - g.lastMouseY

			g.offsetX += dx
			g.offsetY += dy

		} else {
			g.mousePressed = true
		}
		g.lastMouseX = x
		g.lastMouseY = y
	} else {
		g.mousePressed = false
	}

	g.offsetX = clamp(g.offsetX, -g.width*int(g.zoom)+g.width, 0)
	g.offsetY = clamp(g.offsetY, -g.height*int(g.zoom)+g.height, 0)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		adjustedX := (x - g.offsetX) / int(g.zoom)
		adjustedY := (y - g.offsetY) / int(g.zoom)

		g.board.Cells[adjustedY][adjustedX].Alive = !g.board.Cells[adjustedY][adjustedX].Alive
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}

	if g.accumulatedTime > g.gameSpeed && !g.paused {
		g.accumulatedTime -= g.gameSpeed
		g.board.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scene := ebiten.NewImage(g.width, g.height)
	scene.Fill(color.White)

	for y := 0; y < g.board.Height; y++ {
		for x := 0; x < g.board.Width; x++ {
			if g.board.Cells[y][x].Alive {
				scene.Set(x, y, color.Black)
			} else {
				scene.Set(x, y, color.White)
			}
		}
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(g.zoom, g.zoom)
	opts.GeoM.Translate(float64(g.offsetX), float64(g.offsetY))

	screen.DrawImage(scene, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func clamp(value, minVal, maxVal int) int {
	return max(minVal, min(maxVal, value))
}

func main() {
	//	height := 30
	//	width := 30
	//	board := NewBoard(width, height)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Conways Game of life")

	game := NewGame(640, 480)

	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
	//	board.Cells[0][1].Alive = true
	//    board.Cells[1][2].Alive = true
	//    board.Cells[2][0].Alive = true
	//    board.Cells[2][1].Alive = true
	//    board.Cells[2][2].Alive = true
	//
	//
	//    board.Cells[3][15].Alive = true
	//    board.Cells[3][16].Alive = true
	//    board.Cells[3][17].Alive = true
	//
	//    for{
	//        board.Update()
	//        cmd := exec.Command("cmd", "/c", "cls")
	//        cmd.Stdout = os.Stdout
	//        cmd.Run()
	//        board.debugPrint()
	//        time.Sleep(100 * time.Millisecond)
	//    }
}
