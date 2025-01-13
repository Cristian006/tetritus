package main

// TetrominoType represents different types of tetrominos
type TetrominoType int

const (
	I TetrominoType = iota
	O
	T
	S
	Z
	J
	L
)

// Tetromino shapes defined as 4x4 matrices
var TetrominoShapes = map[TetrominoType][][]bool{
	I: {
		{false, false, false, false},
		{true, true, true, true},
		{false, false, false, false},
		{false, false, false, false},
	},
	O: {
		{false, true, true, false},
		{false, true, true, false},
		{false, false, false, false},
		{false, false, false, false},
	},
	T: {
		{false, true, false, false},
		{true, true, true, false},
		{false, false, false, false},
		{false, false, false, false},
	},
	S: {
		{false, true, true, false},
		{true, true, false, false},
		{false, false, false, false},
		{false, false, false, false},
	},
	Z: {
		{true, true, false, false},
		{false, true, true, false},
		{false, false, false, false},
		{false, false, false, false},
	},
	J: {
		{true, false, false, false},
		{true, true, true, false},
		{false, false, false, false},
		{false, false, false, false},
	},
	L: {
		{false, false, true, false},
		{true, true, true, false},
		{false, false, false, false},
		{false, false, false, false},
	},
}

// Colors for different tetrominos
var TetrominoColors = map[TetrominoType]string{
	I: "🟦", // Cyan
	O: "🟨", // Yellow
	T: "🟪", // Purple
	S: "🟩", // Green
	Z: "🟥", // Red
	J: "🟦", // Blue
	L: "🟧", // Orange
}
