package main

import (
	"fmt"
	"log"
	"os"

	"./ui"
	"golang.org/x/crypto/ssh/terminal"
)

func testFunc() {
}

func testMenu(m *ui.Menu) {
	var sLeft ui.Submenu
	sLeft.Add("Rescan", testFunc)
	m.AddWithSubmenu("Left", &sLeft)

	var sFile ui.Submenu
	sFile.Add("Exit", testFunc)
	m.AddWithSubmenu("File", &sFile)

	m.Add("Command")
	m.Add("Options")
	m.Add("Right")

	//m.OpenMenu(0)
}

func main() {
	var localui ui.UI
	var mainMenu ui.Menu
	var leftPanel ui.FilePanel
	var rightPanel ui.FilePanel

	leftPanel.GoTo("/home/artur")
	rightPanel.GoTo("/home/artur")
	testMenu(&mainMenu)
	localui.Init()
	localui.AddMenu(&mainMenu)
	localui.AddFilePanel(&leftPanel)
	localui.AddS3Panel(&rightPanel)

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
		fmt.Println("\rKeys pressed:", stdin)

	}

	//var msgbox ui.MsgBox
	//msgbox.Draw("test")
}
