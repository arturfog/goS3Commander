package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/arturfog/colors"
	"golang.org/x/crypto/ssh/terminal"
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
		ui.rightPanel.Draw(ui.leftPanel.maxWidth, 2)
	}

	if ui.menu != nil {
		ui.menu.DrawMenu()
	}
}

func (ui *UI) GetTerminalSize() (width int, height int) {
	w, h, _ := terminal.GetSize(0)
	return w, h
}

// -------------------- BOTTOM MENU END ---------------------- //

type BottomMenu struct {
	actions []func()
	items   []string
}

func (bm *BottomMenu) getItems() []string {
	return bm.items
}

// -------------------- BOTTOM MENU END ---------------------- //

type Panel struct {
	items  []string
	size   []int
	modfiy []int

	width       int
	maxWidth    int
	X           int
	Y           int
	selectedIdx int
	active      bool

	terminal_w int
	terminal_h int
}

func (p *Panel) Add(name string) {
	p.items = append(p.items, name)
}

func (p *Panel) clear() {
	p.items = nil
}

func (p *Panel) Down() {
	p.selectedIdx += 1
}

func (p *Panel) SetActive(state bool) {
	p.active = state
}

func (p *Panel) Active() bool {
	return p.active
}

func (p *Panel) Up() {
	p.selectedIdx -= 1
}

func (p *Panel) getItems() []string {
	return p.items
}

// -------------------- PANEL END ---------------------- //

type FilePanel struct {
	Panel
	location     string
	location_arr []string
}

func (fp *FilePanel) GoUp() {
	loc := strings.Join(fp.location_arr[:len(fp.location_arr)-1], "/")
	fp.GoTo(loc)
}

func (fp *FilePanel) GoTo(location string) {
	fp.location = location
	fp.location_arr = strings.Split(fp.location, "/")
	fp.clear()
	files, err := ioutil.ReadDir(fp.location)
	if err != nil {
		log.Fatal(err)
	}

	fp.Add("..")
	for _, f := range files {
		fp.Add(f.Name())
	}
	fp.terminal_w, fp.terminal_h, err = terminal.GetSize(0)
	if err == nil {
		fp.maxWidth = fp.terminal_w / 2
	} else {
		fp.maxWidth = 30
	}
}

func (fp *FilePanel) Action() {
	fp.GoTo(fmt.Sprintf("%s/%s", fp.location, fp.items[fp.selectedIdx]))
	fp.selectedIdx = 0
}

func (fp *FilePanel) Draw(X int, Y int) {
	fp.X = X
	fp.Y = Y
	if fp.X > fp.terminal_w {
		fp.X = fp.terminal_w
	}
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
	for idx, element := range fp.getItems() {
		fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y, fp.X)
		fp.Y += 1
		if idx == fp.selectedIdx && fp.active {
			fmt.Printf("\033[0;%d;39m\u2502 %s", colors.BgCyan, element)
		} else {
			fmt.Printf("\033[0;%d;39m\u2502 %s", colors.BgBlue, element)
		}
		for i := 0; i < (fp.maxWidth - len(element) - 4); i++ {
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
	fmt.Printf("\x1b7\x1b[%d;%dH", fp.Y, fp.X)
	// restore cursor position
	fmt.Println("\x1b8")
	fmt.Println("\033[0m\r")
}

// -------------------- FILEPANEL END ---------------------- //
