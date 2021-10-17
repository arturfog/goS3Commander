package main

import (
	"log"

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

	m.OpenMenu(0)
}

func main() {
	var localui ui.UI
	var mainMenu ui.Menu
	var leftPanel ui.FilePanel
	leftPanel.GoTo("/home/artur")
	testMenu(&mainMenu)
	localui.Init()
	localui.AddMenu(&mainMenu)
	localui.AddFilePanel(&leftPanel)

	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}
	defer func() {
		if err := terminal.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()
	/*
		var keys []byte
		in := bufio.NewReader(os.Stdin)
		for {
			_, err := in.Read(keys)
			if err != nil {
				log.Println("stdin:", err)
				break
			}
			fmt.Printf("read rune %q\r\n", keys)
			if keys[0] == 'q' {
				break
			}
		}
	*/
	localui.Redraw()
}
