package ui

import (
	"fmt"
	"io/ioutil"
	"log"

	"../colors"
)

type UI struct {
	menu       *Menu
	leftPanel  *FilePanel
	rightPanel *FilePanel
}

func (ui *UI) Init() {
	ui.menu = nil
	ui.leftPanel = nil
}

func (ui *UI) AddMenu(m *Menu) {
	ui.menu = m
}

func (ui *UI) AddFilePanel(fp *FilePanel) {
	ui.leftPanel = fp
}

func (ui *UI) AddS3Panel(fp *FilePanel) {
	ui.rightPanel = fp
}

func (ui *UI) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (ui *UI) Redraw() {
	ui.clearScreen()

	if ui.leftPanel != nil {
		ui.leftPanel.Draw(0, 2)
	}

	if ui.rightPanel != nil {
		ui.rightPanel.Draw(ui.leftPanel.maxWidth+5, 2)
	}

	if ui.menu != nil {
		DrawMenu(ui.menu)
	}
}

func (ui *UI) GetTerminalSize() (width int, height int) {
	return 80, 25
}

// -------------------- MENU START ---------------------- //

type Menu struct {
	items   []string
	openIdx int
	isOpen  bool
	s       []Submenu
}

func (m *Menu) Add(name string) {
	m.items = append(m.items, name)
	m.openIdx = -1
	m.s = append(m.s, Submenu{})
}

func (m *Menu) AddWithSubmenu(name string, s *Submenu) {
	m.items = append(m.items, name)
	m.openIdx = -1
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

func (m *Menu) OpenMenu(idx int) {
	if idx == len(m.items) {
		idx = 0
	}
	if idx >= 0 && idx < len(m.items) {
		m.openIdx = idx
	}
}

func (m *Menu) CloseMenu() {
	m.openIdx = -1
}

func DrawMenu(m *Menu) {
	fmt.Print("\x1b7\x1b[1;1H")
	for idx, element := range m.getItems() {
		if idx == m.openIdx {
			fmt.Printf("\033[%d;%dm %s    ", colors.BgBlack, colors.FgWhite, element)
		} else {
			fmt.Printf("\033[%d;%dm %s    ", colors.BgCyan, colors.FgBlack, element)
		}
	}
	fmt.Println("\033[0m\r")
	if m.openIdx >= 0 {
		if len(m.items) >= m.openIdx {
			if len(m.s[m.openIdx].items) > 0 {
				DrawSubmenu(&m.s[m.openIdx], 5)
			}
		}
	}
	fmt.Println("\x1b8")
}

// -------------------- MENU END ---------------------- //

type BottomMenu struct {
	actions []func()
	items   []string
}

func (bm *BottomMenu) getItems() []string {
	return bm.items
}

// -------------------- BOTTOM MENU END ---------------------- //

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

type Panel struct {
	items    []string
	width    int
	maxWidth int
	X        int
	Y        int
}

func (p *Panel) Add(name string) {
	p.items = append(p.items, name)
}

func (p *Panel) clear() {
	p.items = nil
}

func (p *Panel) getItems() []string {
	return p.items
}

// -------------------- PANEL END ---------------------- //

type FilePanel struct {
	Panel
	location string
}

func (fp *FilePanel) GoTo(location string) {
	fp.location = location
	fp.clear()
	files, err := ioutil.ReadDir(fp.location)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fp.Add(f.Name())
		if len(f.Name()) > fp.maxWidth {
			fp.maxWidth = len(f.Name())
		}
	}

	fp.maxWidth = fp.maxWidth + 20
}

func (fp *FilePanel) Draw(X int, Y int) {
	fp.X = X
	fp.Y = Y
	fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y, fp.X)
	// left corner
	fmt.Printf("\033[%d;39m\u250c", colors.BgBlue)
	for i := 0; i < fp.maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2510\r")

	fmt.Printf("\x1b[%d;%dH", fp.Y+1, fp.X)
	fmt.Println("\033[44;39m\u2502 \033[1;93mName\033[39m \u2502 \033[1;93mSize \033[39m\u2502 \033[1;93mModify time \033[39m\u2502\r")

	fmt.Printf("\x1b[%d;%dH", fp.Y+2, fp.X)
	fp.Y += 2
	for _, element := range fp.getItems() {
		fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y, fp.X)
		fp.Y += 1
		fmt.Printf("\033[0;%d;39m\u2502 %s", colors.BgBlue, element)
		for i := 0; i < (fp.maxWidth - len(element)); i++ {
			if i+len(element) == 25 {
				fmt.Print(" \u2502")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println(" \u2502\r")

	}
	fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y, fp.X)
	// left corner
	fmt.Printf("\u2514")
	for i := 0; i < fp.maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2518\r")
	fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y+4, fp.X)
	// restore cursor position
	fmt.Println("\x1b8")
	fmt.Println("\033[0m\r")
}

// -------------------- FILEPANEL END ---------------------- //

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
