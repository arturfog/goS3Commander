package ui

import (
	"fmt"

	"github.com/arturfog/colors"
	"golang.org/x/crypto/ssh/terminal"
)

type BottomMenu struct {
	actions []func()
	items   []string

	term_w int
	term_h int
}

func (bm *BottomMenu) Add(name string) {
	bm.items = append(bm.items, name)

	bm.term_w, bm.term_h, _ = terminal.GetSize(0)
}

func (bm *BottomMenu) Draw() {
	fmt.Printf("\x1b7\x1b[%d;1H", bm.term_h-3)
	for idx, element := range bm.getItems() {
		fmt.Printf("\033[%d;%dm %d %s    ", colors.BgCyan, colors.FgBlack, idx+1, element)
	}
	fmt.Println("\x1b8")
}

func (bm *BottomMenu) getItems() []string {
	return bm.items
}
