package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yairp7/go-pathfinding/pf"
	"github.com/yairp7/go-pathfinding/utils"
)

//go:embed map.dat
var mapdata []byte

var SCREEN_SIZE utils.Vec2 = utils.NewVec2(640, 480)
var TILE_SIZE utils.Vec2 = utils.NewVec2(32, 32)

const (
	SetStartState int = iota
	SetEndState   int = iota
)

type Engine struct {
	screenSize utils.Vec2
	pfMap      pf.PFMap
	start, end utils.Vec2
	state      int
	path       []*pf.PFTile
}

func NewEngine(args ...string) *Engine {
	var screenSize utils.Vec2
	var tileSize utils.Vec2

	if len(args) < 3 {
		screenSize = SCREEN_SIZE
		tileSize = TILE_SIZE
	} else {
		w, _ := strconv.Atoi(args[0])
		h, _ := strconv.Atoi(args[1])
		s, _ := strconv.Atoi(args[2])
		screenSize = utils.NewVec2(float64(w), float64(h))
		tileSize = utils.NewVec2(float64(s), float64(s))
	}

	e := &Engine{
		screenSize: screenSize,
	}

	ebiten.SetWindowClosingHandled(true)

	e.pfMap = pf.NewPFMap(screenSize, tileSize)
	e.loadMap()

	return e
}

func (e *Engine) loadMap() {
	var data [][]int
	err := json.Unmarshal(mapdata, &data)
	if err != nil {
		panic(err)
	}

	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[0]); col++ {
			if data[row][col] == 1 {
				e.pfMap.SetObstacle(row, col)
			}
		}
	}
}

func (e *Engine) onMouseClicked() {
	x, y := ebiten.CursorPosition()
	p := utils.NewVec2(float64(x), float64(y))

	if !e.pfMap.IsTileAvailable(x, y) {
		return
	}

	if e.state == SetStartState {
		if e.path != nil {
			e.pfMap.ClearPath(e.path)
		}
		e.pfMap.ClearTile(int(e.start.X), int(e.start.Y))
		e.pfMap.ClearTile(int(e.end.X), int(e.end.Y))
		e.pfMap.SetStartPoint(int(p.X), int(p.Y))
		e.start = p
		e.state = SetEndState
	} else if e.state == SetEndState {
		e.pfMap.SetEndPoint(int(p.X), int(p.Y))
		e.end = p
		e.state = SetStartState
		e.calcPath()
	}
}

func (e *Engine) calcPath() {
	data := e.pfMap.ShortestPath(int(e.start.X), int(e.start.Y), int(e.end.X), int(e.end.Y))

	path := utils.NewStack[*pf.PFTile]()
	node := &data
	path.Push(node.Tile)
	for node.Parent != nil {
		path.Push(node.Parent.Tile)
		node = node.Parent
	}

	e.path = make([]*pf.PFTile, path.Size())
	i := 0
	for path.Size() > 0 {
		e.path[i] = path.Pop()
		i++
	}

	e.pfMap.SetPath(e.path)
}

func (e *Engine) Update() error {
	if ebiten.IsWindowBeingClosed() {
		return ebiten.Termination
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		e.onMouseClicked()
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.pfMap.Draw(screen)
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(e.screenSize.X), int(e.screenSize.Y)
}

func Start(args ...string) {
	ebiten.SetWindowSize(int(SCREEN_SIZE.X), int(SCREEN_SIZE.Y))
	ebiten.SetWindowTitle("Go Pathfinding (A*)")
	if err := ebiten.RunGame(NewEngine(args...)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	argsWithoutProg := os.Args
	Start(argsWithoutProg...)
}
