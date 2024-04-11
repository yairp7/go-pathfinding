package pf

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yairp7/go-pathfinding/utils"
)

const (
	TileStateUnknown  int = iota
	TileStateMarked       = iota
	TileStatePath         = iota
	TileStateStart        = iota
	TileStateEnd          = iota
	TileStateObstacle     = iota
)

var (
	TileStateStartColor    color.RGBA = color.RGBA{R: 255, G: 255, B: 0, A: 255}
	TileStateEndColor      color.RGBA = color.RGBA{R: 0, G: 255, B: 255, A: 255}
	TileStatePathColor     color.RGBA = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	TileStateObstacleColor color.RGBA = color.RGBA{R: 127, G: 50, B: 127, A: 255}
)

type SearchNodeData struct {
	Tile   *PFTile
	Parent *SearchNodeData
}

type PFTile struct {
	T        int
	Row, Col int
	P        utils.Vec2
}

type PFMap struct {
	tiles         [][]PFTile
	mapSize       utils.Vec2
	tileSize      utils.Vec2
	graph         navGraph[*PFTile]
	numGraphNodes int
}

func NewPFMap(mapSize utils.Vec2, tileSize utils.Vec2) PFMap {
	m := PFMap{
		mapSize:  mapSize,
		tileSize: tileSize,
	}
	m.tiles = createTiles(mapSize, tileSize)
	m.graph = m.bfsBuildGraph()
	return m
}

func (m PFMap) ShortestPath(
	fromX, fromY int,
	toX, toY int,
) SearchNodeData {
	fromTile := m.tileByCoords(fromX, fromY)
	toTile := m.tileByCoords(toX, toY)
	return m.aStarAlgo(
		fromTile.Row, fromTile.Col,
		toTile.Row, toTile.Col,
		func(current, target *PFTile) float64 {
			return math.Abs(current.P.X-target.P.X) + math.Abs(current.P.Y-target.P.Y)
		},
	)
}

func (m PFMap) tileToColor(tileType int) color.RGBA {
	switch tileType {
	case TileStateStart:
		return TileStateStartColor
	case TileStateEnd:
		return TileStateEndColor
	case TileStatePath:
		return TileStatePathColor
	case TileStateObstacle:
		return TileStateObstacleColor
	}
	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

func (m PFMap) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	gap := m.tileSize.Scale(0.05)
	drawTileSize := m.tileSize.Scale(0.9)
	for row := 0; row < len(m.tiles); row++ {
		for col := 0; col < len(m.tiles[0]); col++ {
			tile := m.tiles[row][col]
			x := float64(col*int(m.tileSize.X)) + gap.X
			y := float64(row*int(m.tileSize.Y)) + gap.Y
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(drawTileSize.X), float32(drawTileSize.Y), m.tileToColor(tile.T), true)
		}
	}
}

func (m PFMap) tileByCoords(x, y int) *PFTile {
	return &m.tiles[y/int(m.tileSize.Y)][x/int(m.tileSize.X)]
}

func (m PFMap) isCoordsValid(x, y int) bool {
	if (x < 0 || x > int(m.mapSize.X)) || (y < 0 || y > int(m.mapSize.Y)) {
		return false
	}
	return true
}

func (m PFMap) isRowColValid(row, col int) bool {
	if (row < 0 || row >= len(m.tiles)) || (col < 0 || col >= len(m.tiles[0])) {
		return false
	}
	return true
}

func (m PFMap) IsTileAvailable(x, y int) bool {
	if !m.isCoordsValid(x, y) {
		return false
	}

	tile := m.tileByCoords(x, y)
	return tile.T == TileStateUnknown
}

func (m *PFMap) ClearTile(x, y int) {
	if !m.isCoordsValid(x, y) {
		return
	}

	tile := m.tileByCoords(x, y)
	tile.T = TileStateUnknown
}

func (m *PFMap) SetStartPoint(x, y int) {
	if !m.isCoordsValid(x, y) {
		return
	}

	tile := m.tileByCoords(x, y)
	tile.T = TileStateStart
}

func (m *PFMap) SetEndPoint(x, y int) {
	if !m.isCoordsValid(x, y) {
		return
	}

	tile := m.tileByCoords(x, y)
	tile.T = TileStateEnd
}

func (m *PFMap) ClearPath(path []*PFTile) {
	for i := 0; i < len(path); i++ {
		path[i].T = TileStateUnknown
	}
}

func (m *PFMap) SetPath(path []*PFTile) {
	for i := 0; i < len(path); i++ {
		path[i].T = TileStatePath
	}
}

func (m *PFMap) SetObstacle(row, col int) {
	if !m.isRowColValid(row, col) {
		return
	}

	m.tiles[row][col].T = TileStateObstacle
}

func createTiles(mapSize utils.Vec2, tileSize utils.Vec2) [][]PFTile {
	columns := int(mapSize.X / tileSize.X)
	rows := int(mapSize.Y / tileSize.Y)

	halfTileSize := tileSize.Scale(0.5)

	tiles := make([][]PFTile, rows)
	for r := 0; r < rows; r++ {
		tiles[r] = make([]PFTile, columns)
		for c := 0; c < columns; c++ {
			tiles[r][c].T = TileStateUnknown
			tiles[r][c].Row = r
			tiles[r][c].Col = c
			tiles[r][c].P = utils.NewVec2(float64(c*int(tileSize.X))+halfTileSize.X, float64(r*int(tileSize.Y))+halfTileSize.Y)
		}
	}

	return tiles
}

func isValidTile(tiles [][]PFTile, row, col int) bool {
	return tiles[row][col].T == TileStateUnknown || tiles[row][col].T == TileStateStart || tiles[row][col].T == TileStateEnd
}

func (s PFMap) bfsBuildGraph() (tilesGraph navGraph[*PFTile]) {
	row, col := 0, 0
	count := 1
	h := len(s.tiles)
	if h == 0 {
		return
	}
	l := len(s.tiles[0])

	visited := make([][]bool, h)
	for i := range visited {
		visited[i] = make([]bool, l)
	}

	canMove := func(row, col int) bool {
		if row < 0 || col < 0 || row >= h || col >= l || visited[row][col] {
			return false
		}

		if !isValidTile(s.tiles, row, col) {
			return false
		}

		return true
	}

	addNode := func(
		origin *PFTile,
		row, col int,
		weight float64,
		queue *utils.Queue[*PFTile],
	) {
		if !canMove(row, col) {
			return
		}

		queue.Push(&s.tiles[row][col])
		tilesGraph.addEdge(origin, &s.tiles[row][col], weight)
		count++
	}

	tilesGraph = newNavigationGraph[*PFTile]()
	queue := utils.NewQueue[*PFTile]()
	queue.Push(&s.tiles[row][col])

	for queue.Size() > 0 {
		p := queue.Pop()

		row, col := p.Row, p.Col

		if !canMove(row, col) {
			continue
		}

		visited[row][col] = true

		addNode(&s.tiles[row][col], row, col-1, 1, queue)
		addNode(&s.tiles[row][col], row, col+1, 1, queue)
		// addNode(&s.tiles[row][col], row-1, col-1, 1.41, queue)
		// addNode(&s.tiles[row][col], row-1, col+1, 1.41, queue)
		addNode(&s.tiles[row][col], row-1, col, 1, queue)
		addNode(&s.tiles[row][col], row+1, col, 1, queue)
		// addNode(&s.tiles[row][col], row+1, col-1, 1.41, queue)
		// addNode(&s.tiles[row][col], row+1, col+1, 1.41, queue)
	}

	return
}

func (s PFMap) aStarAlgo(
	fromRow, fromCol int,
	toRow, toCol int,
	hFunc func(current, target *PFTile) float64,
) SearchNodeData {
	origin := &s.tiles[fromRow][fromCol]
	destination := &s.tiles[toRow][toCol]

	h := utils.NewMinHeap[SearchNodeData]()
	h.Push(
		utils.HeapNode[SearchNodeData]{
			TotalWeight: 0,
			Data:        SearchNodeData{Tile: origin, Parent: nil},
		},
	)
	visited := make(map[*PFTile]bool)

	for len(*h.Values) > 0 {
		p := h.Pop()
		node := p.Data.Tile

		if visited[node] || !isValidTile(s.tiles, node.Row, node.Col) {
			continue
		}

		if node == destination {
			return p.Data
		}

		for _, e := range s.graph.getEdges(node) {
			if visited[e.node] {
				continue
			}

			h.Push(
				utils.HeapNode[SearchNodeData]{
					TotalWeight: p.TotalWeight + e.weight + hFunc(e.node, destination),
					Data:        SearchNodeData{Tile: e.node, Parent: &p.Data},
				},
			)
		}

		visited[node] = true
	}

	return SearchNodeData{}
}
