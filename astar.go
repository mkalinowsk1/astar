package main

import (
	"astar/colors"
	"container/heap"
	"math"
	"os"
	"reflect"

	"github.com/veandco/go-sdl2/sdl"
)

type Spot struct {
	row        int32
	col        int32
	width      int32
	total_rows int32
	color      colors.Color
	neighbours []*Spot
}

func (s *Spot) GetPos() (int32, int32) {
	return s.col, s.row
}

func (s *Spot) IsClosed() bool {
	return reflect.DeepEqual(s.color, colors.Red())
}

func (s *Spot) IsOpen() bool {
	return reflect.DeepEqual(s.color, colors.Green())
}

func (s *Spot) IsBarrier() bool {
	return reflect.DeepEqual(s.color, colors.Black())
}

func (s *Spot) IsStart() bool {
	return reflect.DeepEqual(s.color, colors.Orange())
}

func (s *Spot) IsEnd() bool {
	return reflect.DeepEqual(s.color, colors.Turquoise())
}

func (s *Spot) Reset() {
	s.color = colors.White()
}

func (s *Spot) MakeStart() {
	s.color = colors.Orange()
}

func (s *Spot) MakeEnd() {
	s.color = colors.Turquoise()
}

func (s *Spot) MakeClosed() {
	s.color = colors.Red()
}

func (s *Spot) MakeOpen() {
	s.color = colors.Green()
}

func (s *Spot) MakeBarrier() {
	s.color = colors.Black()
}

func (s *Spot) MakePath() {
	s.color = colors.Purple()
}

func (s *Spot) Draw(renderer *sdl.Renderer) {
	rect := sdl.Rect{X: s.row * s.width, Y: s.col * s.width, W: s.width, H: s.width}
	renderer.SetDrawColor(s.color.RGBA())
	renderer.FillRect(&rect)
}

func (s *Spot) UpdateNeighbours(grid *[][]Spot) {
	s.neighbours = make([]*Spot, 0, 4)
	if s.row < (s.total_rows-1) && !((*grid)[s.row+1][s.col].IsBarrier()) {
		s.neighbours = append(s.neighbours, &(*grid)[s.row+1][s.col])
	}
	if s.row > 0 && !((*grid)[s.row-1][s.col].IsBarrier()) {
		s.neighbours = append(s.neighbours, &(*grid)[s.row-1][s.col])
	}
	if s.col < (s.total_rows-1) && !((*grid)[s.row][s.col+1].IsBarrier()) {
		s.neighbours = append(s.neighbours, &(*grid)[s.row][s.col+1])
	}
	if s.col > 0 && !((*grid)[s.row][s.col-1].IsBarrier()) {
		s.neighbours = append(s.neighbours, &(*grid)[s.row][s.col-1])
	}
}

type PriorityQueue []*Item

type Item struct {
	value    *Spot
	priority int
	index    int
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Update(item *Item, value *Spot, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (pq PriorityQueue) Contains(spot *Spot) bool {
	for _, item := range pq {
		if item.value == spot {
			return true
		}
	}
	return false
}

func h(x1, y1, x2, y2 int32) int {
	return int(math.Abs(float64(x1-x2))) + int(math.Abs(float64(y1-y2)))

}

func makeGrid(rows, width int32) [][]Spot {
	grid := make([][]Spot, rows)
	gap := width / rows
	var i, j int32
	for i = 0; i < rows; i++ {
		grid[i] = make([]Spot, rows)
		for j = 0; j < rows; j++ {
			spot := Spot{i, j, gap, rows, colors.White(), make([]*Spot, 0, 4)}
			grid[i][j] = spot
		}

	}
	return grid
}

func drawGrid(renderer *sdl.Renderer, rows, width int32) {
	gap := width / rows
	var i, j int32
	for i = 0; i < rows; i++ {
		renderer.SetDrawColor(colors.Grey().RGBA())
		renderer.DrawLine(0, i*gap, width, i*gap)
		for j = 0; j < rows; j++ {
			renderer.SetDrawColor(colors.Grey().RGBA())
			renderer.DrawLine(j*gap, 0, j*gap, width)
		}
	}
}
func draw(renderer *sdl.Renderer, grid [][]Spot, rows, width int32) {
	renderer.SetDrawColor(colors.White().RGBA())
	renderer.Clear()

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			grid[i][j].Draw(renderer)
		}
	}

	drawGrid(renderer, rows, width)
	renderer.Present()
	sdl.Delay(16)

}

func reconstructPath(renderer *sdl.Renderer, grid [][]Spot, rows, width int32, came_from map[*Spot]*Spot, current *Spot) {
	for current != nil {
		current.MakePath()
		current = came_from[current]
		draw(renderer, grid, rows, width)
	}
}
func algorithm(renderer *sdl.Renderer, grid [][]Spot, rows, width int32, start *Spot, end *Spot) bool {

	openSet := make(PriorityQueue, 0)
	heap.Init(&openSet)
	heap.Push(&openSet, &Item{value: start, priority: 0})
	cameFrom := make(map[*Spot]*Spot)
	gScore := make(map[*Spot]int)
	fScore := make(map[*Spot]int)

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			gScore[&grid[i][j]] = math.MaxInt32
			fScore[&grid[i][j]] = math.MaxInt32
		}
	}
	gScore[start] = 0
	fScore[start] = h(start.col, start.row, end.col, end.row)

	for openSet.Len() > 0 {
		current := heap.Pop(&openSet).(*Item).value

		if current == end {
			reconstructPath(renderer, grid, rows, width, cameFrom, end)
			end.MakeEnd()
			return true
		}

		for _, neighbour := range current.neighbours {
			tempGScore := gScore[current] + 1

			if tempGScore < gScore[neighbour] {
				cameFrom[neighbour] = current
				gScore[neighbour] = tempGScore
				fScore[neighbour] = (tempGScore + h(neighbour.col, neighbour.row, end.col, end.row))

				if !openSet.Contains(neighbour) {
					heap.Push(&openSet, &Item{value: neighbour, priority: fScore[neighbour]})
					neighbour.MakeOpen()
				}

			}

		}
		draw(renderer, grid, rows, width)

		if current != start {
			current.MakeClosed()
		}
	}

	return false
}

func findPath(renderer *sdl.Renderer, grid [][]Spot, rows, width int32, start, end *Spot, done chan bool) {
	algorithm(renderer, grid, rows, width, start, end)

	done <- true
}

func getClickedPos(x, y, rows, width int32) (int32, int32) {
	gap := width / rows

	row := y / gap
	col := x / gap

	return row, col

}

func run(width int32) (er error) {
	var window *sdl.Window
	var renderer *sdl.Renderer

	window, err := sdl.CreateWindow("astar", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		600, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	var ROWS int32 = 50
	grid := makeGrid(ROWS, width)

	var start *Spot
	var zeros *Spot
	var end *Spot
	var zeroe *Spot
	done := make(chan bool)

	running := true
	for running {
		draw(renderer, grid, ROWS, width)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch event.GetType() {
			case sdl.QUIT:
				running = false
			case sdl.MOUSEBUTTONDOWN:
				mouseEvent := event.(*sdl.MouseButtonEvent)
				if mouseEvent.Button == sdl.BUTTON_LEFT {
					x, y, _ := sdl.GetMouseState()
					row, col := getClickedPos(y, x, ROWS, width)
					spot := &grid[row][col]
					if start == nil && spot != end {
						start = spot
						start.MakeStart()
					} else if end == nil && spot != start {
						end = spot
						end.MakeEnd()
					} else if spot != end && spot != start {
						spot.MakeBarrier()
					}
				} else if mouseEvent.Button == sdl.BUTTON_RIGHT {
					x, y, _ := sdl.GetMouseState()
					row, col := getClickedPos(y, x, ROWS, width)
					spot := &grid[row][col]
					spot.Reset()
					if spot == start {
						start = zeros
					} else if spot == end {
						end = zeroe
					}

				}
			case sdl.KEYDOWN:
				keyEvent := event.(*sdl.KeyboardEvent)
				if keyEvent.Keysym.Sym == sdl.K_SPACE && start != nil && end != nil {
					for i := 0; i < len(grid); i++ {
						for j := 0; j < len(grid[i]); j++ {
							grid[i][j].UpdateNeighbours(&grid)
						}
					}
					go findPath(renderer, grid, ROWS, width, start, end, done)
					<-done
				}
			}

		}
		renderer.Present()
	}
	sdl.Quit()
	return
}

func main() {

	if er := run(600); er != nil {
		os.Exit(1)
	}
}
