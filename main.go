package main

import (
	"log"
	"os"

	"github.com/arturfog/ui"
	"golang.org/x/crypto/ssh/terminal"
)

func testFunc() {
}

func testMenu(m *ui.Menu) {
	var sLeft ui.Submenu
	sLeft.Add("Sort ...", testFunc)
	sLeft.Add("Filter", testFunc)
	sLeft.Add("S3 link ...", testFunc)
	sLeft.Add("Rescan", testFunc)
	m.AddWithSubmenu("Left", &sLeft)

	var sFile ui.Submenu
	sFile.Add("Exit", testFunc)
	m.AddWithSubmenu("File", &sFile)

	m.Add("Command")
	m.Add("Options")
	m.Add("Right")
}

func main() {
	var localui ui.UI
	var mainMenu ui.Menu
	var leftPanel ui.FilePanel
	var rightPanel ui.FilePanel
	var bottomMenu ui.BottomMenu

	leftPanel.GoTo("/Users/arturfogiel")
	leftPanel.SetActive(true)
	rightPanel.GoTo("/Users/arturfogiel")
	testMenu(&mainMenu)
	localui.Init()
	localui.AddMenu(&mainMenu)
	localui.AddFilePanel(&leftPanel)
	localui.AddS3Panel(&rightPanel)
	bottomMenu.Add("Help")
	bottomMenu.Add("Menu")
	bottomMenu.Add("View")
	bottomMenu.Add("Edit")
	bottomMenu.Add("Copy")
	bottomMenu.Add("Move")
	bottomMenu.Add("Mkdir")
	bottomMenu.Add("Delete")
	bottomMenu.Add("Quit")
	localui.AddBottomMenu(&bottomMenu)

	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}
	defer func() {
		if err := terminal.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	ch := make(chan []byte, 5)
	go func(ch chan []byte) {

		for {
			var b = make([]byte, 5)
			os.Stdin.Read(b)
			ch <- b
		}
	}(ch)

	localui.Redraw()

	for {
		stdin, _ := <-ch
		if string(stdin[0]) == "q" {
			break
		}
		// F9
		if stdin[2] == 50 && stdin[3] == 48 {
			mainMenu.OpenMenu(0)
			localui.Redraw()
		}
		// Escape
		if stdin[0] == 27 && stdin[1] == 0 {
			mainMenu.CloseMenu()
			localui.Redraw()
		}
		// Tab
		if stdin[0] == 9 && stdin[1] == 0 {
			if leftPanel.Active() {
				rightPanel.SetActive(true)
				leftPanel.SetActive(false)
			} else {
				rightPanel.SetActive(false)
				leftPanel.SetActive(true)
			}
			localui.Redraw()
		}
		// Enter
		if stdin[0] == 13 && stdin[1] == 0 {
			if mainMenu.IsOpen() {

			} else {
				if leftPanel.Active() {
					leftPanel.Action()
				} else {
					rightPanel.Action()
				}
			}
			localui.Redraw()
		}
		// Right
		// 27 91 67
		if stdin[0] == 27 && stdin[1] == 91 && stdin[2] == 67 {
			if mainMenu.IsOpen() {
				mainMenu.OpenMenu(mainMenu.GetOpenedIdx() + 1)
			}
			localui.Redraw()
		}
		// Up
		// 27 91 65
		if stdin[0] == 27 && stdin[1] == 91 && stdin[2] == 65 {
			if mainMenu.IsOpen() {
				mainMenu.Up()
			} else {
				if leftPanel.Active() {
					leftPanel.Up()
				} else {
					rightPanel.Up()
				}
			}
			localui.Redraw()
		}
		// Down
		// 27 91 66
		if stdin[0] == 27 && stdin[1] == 91 && stdin[2] == 66 {
			if mainMenu.IsOpen() {
				mainMenu.Down()
			} else {
				if leftPanel.Active() {
					leftPanel.Down()
				} else {
					rightPanel.Down()
				}
			}
			localui.Redraw()
		}
		// Left
		// 27 91 68
		if stdin[0] == 27 && stdin[1] == 91 && stdin[2] == 68 {
			if mainMenu.IsOpen() {
				mainMenu.OpenMenu(mainMenu.GetOpenedIdx() - 1)
			}
			localui.Redraw()
		}
		//fmt.Println("\rKeys pressed:", stdin)
	}
}
