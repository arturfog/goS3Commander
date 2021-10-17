package main

//
import (
	"./ui"
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
	localui.Redraw()
}
