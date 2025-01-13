package main

type GameStyle int

const (
	EmojiStyle GameStyle = iota
	ClassicStyle
)

// Block styles for different game modes
var BlockStyles = map[GameStyle]map[TetrominoType]string{
	EmojiStyle: {
		I: "🟦", // Cyan
		O: "🟨", // Yellow
		T: "🟪", // Purple
		S: "🟩", // Green
		Z: "🟥", // Red
		J: "🟦", // Blue
		L: "🟧", // Orange
	},
	ClassicStyle: {
		I: "[]",
		O: "[]",
		T: "[]",
		S: "[]",
		Z: "[]",
		J: "[]",
		L: "[]",
	},
}

var EmptyBlock = map[GameStyle]string{
	EmojiStyle:   "⬜",
	ClassicStyle: ". ",
}
