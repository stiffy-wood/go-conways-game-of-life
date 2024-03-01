package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
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
            if b.Cells[y][x].Alive{
                symbol = "[#]"
            }
			fmt.Print(symbol, " ")
		}
        fmt.Print("\n")
	}
    fmt.Print("\n")
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
			if b.Cells[curY + y][curX + x].Alive {
				neighbors++
			}
		}
	}
	return neighbors
}

func main() {
	height := 30
	width := 30
	board := NewBoard(width, height)

	board.Cells[0][1].Alive = true
    board.Cells[1][2].Alive = true
    board.Cells[2][0].Alive = true
    board.Cells[2][1].Alive = true
    board.Cells[2][2].Alive = true

    for{
        board.Update()
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
        board.debugPrint()
        time.Sleep(10 * time.Millisecond)
    }
}
