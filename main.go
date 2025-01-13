package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Game represents the main game state
type Game struct {
	board        [][]string
	currentPiece Piece
	score        int
	gameOver     bool
	style        GameStyle
	selecting    bool // true when selecting style
}

// Piece represents a tetris piece
type Piece struct {
	shape     [][]bool
	x, y      int
	pieceType TetrominoType
}

func initialModel() Game {
	return Game{
		selecting: true,
		style:     EmojiStyle,
	}
}

func (g *Game) spawnNewPiece() {
	pieceTypes := []TetrominoType{I, O, T, S, Z, J, L}
	newType := pieceTypes[rand.Intn(len(pieceTypes))]

	// Make a deep copy of the shape
	shape := make([][]bool, 4)
	for i := range shape {
		shape[i] = make([]bool, 4)
		copy(shape[i], TetrominoShapes[newType][i])
	}

	g.currentPiece = Piece{
		shape:     shape,
		x:         3,
		y:         0,
		pieceType: newType,
	}
}

func (g *Game) checkCollision(piece Piece) bool {
	// First check if piece is initialized
	if piece.shape == nil {
		return false
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if piece.shape[y][x] {
				newY := piece.y + y
				newX := piece.x + x
				if newX < 0 || newX >= 10 || newY >= 20 {
					return true
				}
				if newY >= 0 && g.board[newY][newX] != EmptyBlock[g.style] {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) lockPiece() {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if g.currentPiece.shape[y][x] {
				boardY := g.currentPiece.y + y
				boardX := g.currentPiece.x + x
				if boardY >= 0 && boardY < 20 && boardX >= 0 && boardX < 10 {
					g.board[boardY][boardX] = BlockStyles[g.style][g.currentPiece.pieceType]
				}
			}
		}
	}
	g.clearLines()
	g.spawnNewPiece()
	if g.checkCollision(g.currentPiece) {
		g.gameOver = true
	}
}

func (g *Game) clearLines() {
	for y := 19; y >= 0; y-- {
		full := true
		for x := 0; x < 10; x++ {
			if g.board[y][x] == EmptyBlock[g.style] {
				full = false
				break
			}
		}
		if full {
			g.score += 100
			// Move all lines above down
			for y2 := y; y2 > 0; y2-- {
				copy(g.board[y2], g.board[y2-1])
			}
			// Clear top line
			for x := 0; x < 10; x++ {
				g.board[0][x] = EmptyBlock[g.style]
			}
			y++ // Check the same line again
		}
	}
}

func (g *Game) rotatePiece() {
	newShape := make([][]bool, 4)
	for i := range newShape {
		newShape[i] = make([]bool, 4)
	}

	// Rotate clockwise
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			newShape[x][3-y] = g.currentPiece.shape[y][x]
		}
	}

	oldShape := g.currentPiece.shape
	g.currentPiece.shape = newShape
	if g.checkCollision(g.currentPiece) {
		g.currentPiece.shape = oldShape
	}
}

func (g Game) Init() tea.Cmd {
	// Don't start ticking until game starts
	if g.selecting {
		return nil
	}
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if g.selecting {
			switch msg.String() {
			case "1":
				g.style = EmojiStyle
				g.selecting = false
				g.initializeBoard()
				g.spawnNewPiece()
				// Return with a tick command to start the game
				return g, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			case "2":
				g.style = ClassicStyle
				g.selecting = false
				g.initializeBoard()
				g.spawnNewPiece()
				// Return with a tick command to start the game
				return g, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			case "q", "ctrl+c":
				return g, tea.Quit
			}
			return g, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "left":
			g.currentPiece.x--
			if g.checkCollision(g.currentPiece) {
				g.currentPiece.x++
			}
		case "right":
			g.currentPiece.x++
			if g.checkCollision(g.currentPiece) {
				g.currentPiece.x--
			}
		case "down":
			g.currentPiece.y++
			if g.checkCollision(g.currentPiece) {
				g.currentPiece.y--
				g.lockPiece()
			}
		case "up":
			g.rotatePiece()
		}
	case tickMsg:
		if !g.gameOver {
			g.currentPiece.y++
			if g.checkCollision(g.currentPiece) {
				g.currentPiece.y--
				g.lockPiece()
			}
		}
		return g, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return g, nil
}

func (g *Game) initializeBoard() {
	g.board = make([][]string, 20)
	for i := range g.board {
		g.board[i] = make([]string, 10)
		for j := range g.board[i] {
			g.board[i][j] = EmptyBlock[g.style]
		}
	}
}

func (g Game) View() string {
	if g.selecting {
		return `
Select Game Style:

1. Emoji Style (🟦 🟨 🟪 etc.)
2. Classic Style ([] . )

Press 1 or 2 to start...
`
	}

	// Create a temporary board for rendering
	tempBoard := make([][]string, 20)
	for i := range tempBoard {
		tempBoard[i] = make([]string, 10)
		copy(tempBoard[i], g.board[i])
	}

	// Draw current piece
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if g.currentPiece.shape[y][x] {
				boardY := g.currentPiece.y + y
				boardX := g.currentPiece.x + x
				if boardY >= 0 && boardY < 20 && boardX >= 0 && boardX < 10 {
					tempBoard[boardY][boardX] = BlockStyles[g.style][g.currentPiece.pieceType]
				}
			}
		}
	}

	// Build the view string
	var s string
	s += "Tetris\n\n"

	// Add border for classic style
	if g.style == ClassicStyle {
		s += "┌" + strings.Repeat("──", 10) + "┐\n"
	}

	// Draw the board
	for _, row := range tempBoard {
		if g.style == ClassicStyle {
			s += "│"
		}
		for _, cell := range row {
			if cell == "" {
				s += EmptyBlock[g.style]
			} else {
				s += cell
			}
		}
		if g.style == ClassicStyle {
			s += "│"
		}
		s += "\n"
	}

	// Add bottom border for classic style
	if g.style == ClassicStyle {
		s += "└" + strings.Repeat("──", 10) + "┘\n"
	}

	s += fmt.Sprintf("\nScore: %d\n", g.score)
	if g.gameOver {
		if g.style == EmojiStyle {
			s += "\nGame Over!\nPress 'q' to quit\n"
		} else {
			s += "\n*** GAME OVER ***\nPress 'q' to quit\n"
		}
	} else {
		s += "\nControls: ←↑↓→ to move/rotate, 'q' to quit\n"
	}

	return s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return
	}
}
