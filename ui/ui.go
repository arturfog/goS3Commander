package ui

import (
	"fmt"
	"io/ioutil"
	"log"

	"../colors"
)

type UI struct {
	menu      *Menu
	filePanel *FilePanel
}

func (ui *UI) Init() {
	ui.menu = nil
	ui.filePanel = nil
}

func (ui *UI) AddMenu(m *Menu) {
	ui.menu = m
}

func (ui *UI) AddFilePanel(fp *FilePanel) {
	ui.filePanel = fp
}

func (ui *UI) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (ui *UI) Redraw() {
	ui.clearScreen()
	if ui.menu != nil {
		DrawMenu(ui.menu)
	}
	if ui.filePanel != nil {
		DrawLeftPanel(ui.filePanel)
	}
}

func (ui *UI) GetTerminalSize() (width int, height int) {
	return 80, 25
}

// -------------------- MENU START ---------------------- //

type Menu struct {
	items   []string
	openIdx int
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

func (m *Menu) getItems() []string {
	return m.items
}

func (m *Menu) OpenMenu(idx int) {
	if idx >= 0 && idx < len(m.items) {
		m.openIdx = idx
	}
}

func (m *Menu) CloseMenu() {
	m.openIdx = -1
}

func DrawMenu(m *Menu) {
	for _, element := range m.getItems() {
		fmt.Printf("\033[%d;%dm %s    ", colors.BgCyan, colors.FgDefault, element)
	}
	fmt.Println("\033[0m\r")
	if m.openIdx >= 0 {
		if len(m.s) > m.openIdx {
			DrawSubmenu(&m.s[m.openIdx])
		}
	}
}

// -------------------- MENU END ---------------------- //

// -------------------- SUBMENU START ---------------------- //

type Submenu struct {
	actions  []func()
	items    []string
	maxWidth int
}

func (s *Submenu) Add(name string, action func()) {
	if len(name) > s.maxWidth {
		s.maxWidth = len(name)
	}
	s.items = append(s.items, name)
	s.actions = append(s.actions, action)
}

func (s *Submenu) getItems() []string {
	return s.items
}

func DrawSubmenu(s *Submenu) {
	var maxWidth = s.maxWidth + 7

	// left corner
	fmt.Print("\033[46;39m\u250c")
	for i := 0; i < maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2510\r")
	for _, element := range s.getItems() {
		fmt.Print(fmt.Sprintf("\u2502  %s", element))
		for i := 0; i < maxWidth-len(element)-2; i++ {
			fmt.Print(" ")
		}
		fmt.Println("\u2502\r")
	}
	// left corner
	fmt.Print("\u2514")
	for i := 0; i < maxWidth; i++ {
		fmt.Print("\u2500")
	}
	// right corner
	fmt.Println("\u2518\r")
}

// -------------------- SUBMENU END ---------------------- //

type Panel struct {
	items    []string
	maxWidth int
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
}

// -------------------- FILEPANEL END ---------------------- //

func DrawLeftPanel(p *FilePanel) {
	p.maxWidth = p.maxWidth + 20
	fmt.Println("\033[44;39m\u2502 \033[1;93mName\033[39m \u2502 \033[1;93mSize \033[39m\u2502 \033[1;93mModify time \033[39m\u2502\r")
	for _, element := range p.getItems() {
		fmt.Printf("\033[0;44;39m\u2502 %s", element)
		for i := 0; i < (p.maxWidth - len(element)); i++ {
			fmt.Print(" ")
		}
		fmt.Println(" \u2502\r")
	}
	fmt.Println("\033[0m\r")
}

func DrawRightPanel(p *Panel) {
	fmt.Println("\u2502 Name \u2502 Size \u2502 Modify time \u2502")
	for _, element := range p.getItems() {
		fmt.Println(fmt.Sprintf("\u2502 %s \u2502", element))
	}
	fmt.Println("")
}
