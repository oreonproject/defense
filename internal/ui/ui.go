package ui

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type MainWindow struct {
	*widgets.QMainWindow
	app *widgets.QApplication
}

func NewMainWindow(app *widgets.QApplication) *MainWindow {
	window := &MainWindow{
		QMainWindow: widgets.NewQMainWindow(nil, 0),
		app:        app,
	}

	// Set window properties
	window.SetWindowTitle("Oreon Defense")
	window.SetMinimumSize2(800, 600)

	// Create central widget and layout
	centralWidget := widgets.NewQWidget(window, 0)
	window.SetCentralWidget(centralWidget)
	layout := widgets.NewQVBoxLayout2(centralWidget)

	// Add title label
	title := widgets.NewQLabel2("Oreon Defense", nil, 0)
	title.SetAlignment(core.Qt__AlignCenter)
	titleFont := gui.NewQFont2("Arial", 18, 1, false)
	title.SetFont(titleFont)
	layout.AddWidget(title, 0, core.Qt__AlignCenter)

	// Add status label
	statusLabel := widgets.NewQLabel2("Status: Running", nil, 0)
	statusLabel.SetAlignment(core.Qt__AlignCenter)
	layout.AddWidget(statusLabel, 0, core.Qt__AlignCenter)

	// Add stretch to push content to the top
	layout.AddStretch(1)

	// Create menu bar
	menuBar := window.MenuBar()
	fileMenu := menuBar.AddMenu2("&File")

	// Add quit action
	quitAction := widgets.NewQAction2("&Quit", window)
	quitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))
	quitAction.ConnectTriggered(func(checked bool) {
		app.Quit()
	})
	fileMenu.AddActions([]*widgets.QAction{quitAction})

	// Add help menu
	helpMenu := menuBar.AddMenu2("&Help")
	aboutAction := widgets.NewQAction2("&About", window)
	aboutAction.ConnectTriggered(func(checked bool) {
		widgets.QMessageBox_About(window, "About Oreon Defense",
			fmt.Sprintf("Oreon Defense v%s\n\nA security monitoring and defense system.", "0.1.0"),
		)
	})
	helpMenu.AddActions([]*widgets.QAction{aboutAction})

	return window
}
