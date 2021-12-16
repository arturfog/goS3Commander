package ui

import (
	"fmt"

	"github.com/arturfog/colors"
	"golang.org/x/crypto/ssh/terminal"
)

// -------------------- MENU START ---------------------- //

type Menu struct {
	items   []string
	startX  []int
	openIdx int
	isOpen  bool
	s       []Submenu
	term_w  int
	term_h  int
	err     error
}

func (m *Menu) Add(name string) {
	m.items = append(m.items, name)
	m.openIdx = -1
	if len(m.startX) > 0 {
		var idx = len(m.startX) - 1
		m.startX = append(m.startX, m.startX[idx]+len(m.items[idx])+6)
	} else {
		m.startX = append(m.startX, 0)
	}
	m.s = append(m.s, Submenu{})

	m.term_w, m.term_h, _ = terminal.GetSize(0)
}

func (m *Menu) AddWithSubmenu(name string, s *Submenu) {
	m.items = append(m.items, name)
	m.openIdx = -1
	if len(m.startX) > 0 {
		var idx = len(m.startX) - 1
		m.startX = append(m.startX, m.startX[idx]+len(m.items[idx])+6)
	} else {
		m.startX = append(m.startX, 0)
	}
	m.s = append(m.s, *s)
}

func (m *Menu) GetOpenedIdx() int {
	return m.openIdx
}

func (m *Menu) getItems() []string {
	return m.items
}

func (m *Menu) Down() {
	m.s[m.openIdx].Down()
}

func (m *Menu) Up() {
	m.s[m.openIdx].Up()
}

func (m *Menu) IsOpen() bool {
	return m.isOpen
}

func (m *Menu) OpenMenu(idx int) {
	m.isOpen = true
	if idx == len(m.items) {
		idx = 0
	}
	if idx >= 0 && idx < len(m.items) {
		m.openIdx = idx
	}
}

func (m *Menu) CloseMenu() {
	m.isOpen = false
	m.openIdx = -1
}

func (m *Menu) DrawMenu() {
	fmt.Print("\x1b7\x1b[1;1H")
	for idx, element := range m.getItems() {
		if idx == m.openIdx {
			fmt.Printf("\033[%d;%dm %s    ", colors.BgBlack, colors.FgWhite, element)
		} else {
			fmt.Printf("\033[%d;%dm %s    ", colors.BgCyan, colors.FgBlack, element)
		}
	}

	var idx = len(m.startX) - 1
	x := m.startX[idx] + len(m.items[idx])
	for i := 0; i < m.term_w-x-(len(m.items)-4); i++ {
		fmt.Printf(" ")
	}

	fmt.Println("\033[0m\r")
	if m.openIdx >= 0 {
		if len(m.items) >= m.openIdx {
			if len(m.s[m.openIdx].items) > 0 {
				DrawSubmenu(&m.s[m.openIdx], m.startX[m.openIdx])
			}
		}
	}
	fmt.Println("\x1b8")
}

// -------------------- MENU END ---------------------- //

// -------------------- SUBMENU START ---------------------- //

type Submenu struct {
	actions  []func()
	items    []string
	maxWidth int
	openIdx  int
}

func (s *Submenu) Add(name string, action func()) {
	if len(name) > s.maxWidth {
		s.maxWidth = len(name)
	}
	s.items = append(s.items, name)
	s.actions = append(s.actions, action)
	s.openIdx = 0
}

func (s *Submenu) Up() {
	s.openIdx -= 1
}

func (s *Submenu) Down() {
	s.openIdx += 1
}

func (s *Submenu) getItems() []string {
	return s.items
}

func DrawSubmenu(s *Submenu, startX int) {
	var maxWidth = s.maxWidth + 7
	Y := 2
	Z := 0
	fmt.Printf("\x1b7\x1b[%d;%dH", Y, startX)
	// left corner
	fmt.Print("\033[46;39m\u250c")
	for i := 0; i < maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2510\r")
	fmt.Printf("\x1b[%d;%dH", Y+1, startX)
	for idx, element := range s.getItems() {
		fmt.Printf("\033[%d;%dm", colors.BgCyan, colors.FgWhite)
		if idx == s.openIdx {
			fmt.Printf("\u2502\033[%d;%dm %s ", colors.BgBlack, colors.FgWhite, element)
		} else {
			fmt.Printf("\u2502\033[%d;%dm %s ", colors.BgCyan, colors.FgWhite, element)
		}
		for i := 0; i < maxWidth-len(element)-2; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("\033[%d;%dm", colors.BgCyan, colors.FgWhite)
		fmt.Print("\u2502")
		fmt.Printf("\033[%d;%dm  \r", colors.BgBlack, colors.FgWhite)
		fmt.Printf("\x1b[%d;%dH", Y+2+idx, startX)
		Z = Y + 2 + idx
	}

	fmt.Printf("\033[%d;%dm", colors.BgCyan, colors.FgWhite)
	// left corner
	fmt.Print("\u2514")
	for i := 0; i < maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Print("\u2518\033[0m")
	fmt.Printf("\033[%d;%dm  \r", colors.BgBlack, colors.FgWhite)
	fmt.Printf("\x1b[%d;%dH", Z+1, startX+2)
	for i := 0; i < maxWidth+2; i++ {
		fmt.Print(" ")
	}
	fmt.Println("\x1b8")
}

// -------------------- SUBMENU END ---------------------- //
