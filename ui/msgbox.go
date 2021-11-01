package ui

import (
	"fmt"

	"../colors"
)

type MsgBox struct {
	X  int
	Y  int
	W  int
	OK string
}

func (m *MsgBox) Draw(text string) {
	m.X = 15
	m.Y = 10
	m.W = 30
	m.OK = "[ OK ]"

	fmt.Printf("\x1b7\x1b[%d;%dH", m.Y, m.X)
	// left corner
	fmt.Printf("\033[%d;%dm\u250c", colors.FgBlack, colors.BgWhite)
	for i := 0; i < m.W; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2510\r")

	fmt.Printf("\x1b[%d;%dH\u2502", m.Y+1, m.X)
	for i := 0; i < m.W; i++ {
		if i == ((m.W - 2 - len(text)) / 2) {
			fmt.Print(text)
			i += len(text)
		}
		fmt.Print(" ")
	}
	fmt.Println("\u2502")

	// left corner
	fmt.Printf("\x1b[%d;%dH\u2514", m.Y+2, m.X)
	for i := 0; i < m.W; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2518\r")

	// OK button area
	fmt.Printf("\x1b[%d;%dH", m.Y+3, m.X)
	for i := 0; i < m.W; i++ {
		if i == ((m.W - len(m.OK)) / 2) {
			fmt.Printf("%s", m.OK)
			i += len(m.OK) - 2
		}
		fmt.Print(" ")
	}

	// left corner
	fmt.Printf("\x1b[%d;%dH\u2514", m.Y+4, m.X)
	for i := 0; i < m.W; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2518\r")

	// restore cursor position
	fmt.Println("\x1b8")
}

// -------------------- MSGBOX END ---------------------- //
