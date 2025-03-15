package main

import (
	"encoding/json"
	"image/color"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// 20x20 game grid, make cellSize 40 for 10x10
// The ratio of startSpeed/ticksPerSecond is
// how many seconds per move. 480/120 = 4 secs. SLOW
const (
	cellSize       = 20
	startSpeed     = 480 
	ticksPerSecond = 120 // leave this alone
	screenWidth    = 400
	screenHeight   = 400
	gridWidth      = screenWidth / cellSize
	gridHeight     = screenHeight / cellSize
	minSpeed       = 10
	speedBoost     = 5
)

type point struct {
	x, y int
}

type Game struct {
	snake         []point
	direction     point
	nextDirection point
	food          point
	gameOver      bool
	gameRunning   bool // 🚦 Start in a paused state
	tickCount     int
	moveSpeed     int
	mutex         sync.Mutex
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	startX := gridWidth / 2
	startY := gridHeight / 2

	game := &Game{
		snake:         []point{{startX, startY}},
		direction:     point{1, 0}, // Start moving right
		nextDirection: point{1, 0},
		moveSpeed:     startSpeed,
		gameRunning:   false, // 🚦 Game starts paused
	}

	game.spawnFood()
	return game
}

func (g *Game) spawnFood() {
	g.food = point{rand.Intn(gridWidth), rand.Intn(gridHeight)}
}

func (g *Game) Update() error {
	if g.gameOver || !g.gameRunning { // 🚫 Don't update if game isn't running
		return nil
	}

	g.tickCount++
	if g.tickCount >= g.moveSpeed {
		g.tickCount = 0
		g.moveSnake()
	}

	return nil
}

func (g *Game) moveSnake() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.direction = g.nextDirection
	head := g.snake[0]
	newHead := point{head.x + g.direction.x, head.y + g.direction.y}

	if newHead.x < 0 || newHead.x >= gridWidth || newHead.y < 0 || newHead.y >= gridHeight {
		g.gameOver = true
		log.Println("Collision with wall!")
		return
	}

	for _, segment := range g.snake {
		if segment == newHead {
			g.gameOver = true
			log.Println("Collision with self!")
			return
		}
	}

	g.snake = append([]point{newHead}, g.snake...)

	if newHead == g.food {
		log.Println("Food eaten! Increasing speed.")
		g.spawnFood()
		if g.moveSpeed > minSpeed {
			g.moveSpeed -= speedBoost
		}
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, segment := range g.snake {
		ebitenutil.DrawRect(screen, float64(segment.x*cellSize), float64(segment.y*cellSize), cellSize, cellSize, color.RGBA{0, 255, 0, 255})
	}

	ebitenutil.DrawRect(screen, float64(g.food.x*cellSize), float64(g.food.y*cellSize), cellSize, cellSize, color.RGBA{255, 0, 0, 255})

	if g.gameOver {
		text.Draw(screen, "Game Over!", basicfont.Face7x13, screenWidth/3, screenHeight/2, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// =========== HTTP SERVER ===========

// I cut down the info returned to the LLM
// maybe more would help???
type GameStateJSON struct {
	Grid      [][]int `json:"grid"`
	//Snake     []point `json:"snake"`
	//Food      point   `json:"food"`
	Direction string  `json:"direction"`
	GameOver  bool    `json:"game_over"`
	//Speed     int     `json:"speed"`
	//TickCount int     `json:"tick_count"`
	//GameRunning bool  `json:"game_running"` // Added to indicate if game is running
}

func (g *Game) getStateHandler(w http.ResponseWriter, r *http.Request) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	grid := make([][]int, gridHeight)
	for i := range grid {
		grid[i] = make([]int, gridWidth)
	}

	for _, segment := range g.snake {
		grid[segment.y][segment.x] = 1
	}

	grid[g.food.y][g.food.x] = 2

	directionStr := map[point]string{
		{0, -1}: "up",
		{0, 1}:  "down",
		{-1, 0}: "left",
		{1, 0}:  "right",
	}[g.direction]

	state := GameStateJSON{
		Grid:       grid,
		//Snake:      g.snake,
		//Food:       g.food,
		Direction:  directionStr,
		GameOver:   g.gameOver,
		//Speed:      g.moveSpeed,
        //TickCount:  g.tickCount,
		//GameRunning: g.gameRunning, // Return running status
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// Accepts move commands via HTTP
func (g *Game) moveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Move string `json:"move"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	switch request.Move {
	case "up":
		if g.direction.y == 0 {
			g.nextDirection = point{0, -1}
		}
	case "down":
		if g.direction.y == 0 {
			g.nextDirection = point{0, 1}
		}
	case "left":
		if g.direction.x == 0 {
			g.nextDirection = point{-1, 0}
		}
	case "right":
		if g.direction.x == 0 {
			g.nextDirection = point{1, 0}
		}
	default:
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Move accepted"))
}

// Start the game if paused, or reset if game over
func (g *Game) startHandler(w http.ResponseWriter, r *http.Request) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.gameOver {
		// Instead of replacing the object, reset fields manually
		startX, startY := gridWidth/2, gridHeight/2
		g.snake = []point{{startX, startY}}
		g.direction = point{1, 0}
		g.nextDirection = point{1, 0}
		g.moveSpeed = startSpeed
		g.tickCount = 0
		g.gameOver = false
		g.spawnFood()

		log.Println("Game reset after game over.")
	}

	g.gameRunning = true
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Game started!"))
}

func startHTTPServer(game *Game) {
	http.HandleFunc("/state", game.getStateHandler)
	http.HandleFunc("/move", game.moveHandler)
	http.HandleFunc("/start", game.startHandler) // ✅ Handles both start and reset
	log.Println("HTTP server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	game := NewGame()
	go startHTTPServer(game)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game")
	ebiten.SetTPS(ticksPerSecond)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
