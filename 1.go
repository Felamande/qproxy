package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Felamande/go-socks5"
	"github.com/Felamande/qproxy/socks5server"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var verTag string
var verCommitHash string

// func init() {
// 	cmdTag := exec.Command("git", "describe", "--abbrev=0", "--tags")
// 	outTag, err := cmdTag.Output()
// 	if err == nil {
// 		verTag = string(outTag)
// 	}

// 	cmdHash := exec.Command("git", "log", "-1", `--format="%H"`)
// 	outHash, err := cmdHash.Output()
// 	if err == nil {
// 		verCommitHash = string(outHash)
// 	}

// }

type HideWindow struct {
	widgets.QMainWindow
}

func (w *HideWindow) CloseEvent(event gui.QCloseEvent_ITF) {
	w.Hide()
	event.QCloseEvent_PTR().Ignore()
}

func (w *HideWindow) ChangeEvent(event core.QEvent_ITF) {
	if event.QEvent_PTR().Type() == core.QEvent__WindowStateChange && w.IsMinimized() {
		w.Hide()
		event.QEvent_PTR().Ignore()
	}
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	sserver := socks5server.NewSocks5Server()

	var (
		validatorGroup    = widgets.NewQGroupBox2("socks5代理服务器", nil)
		validatorLabel    = widgets.NewQLabel2("端口:", nil, 0)
		validatorLineEdit = widgets.NewQLineEdit(nil)
		StartButton       = widgets.NewQPushButton(nil)
		StopButton        = widgets.NewQPushButton(nil)
		infoLineEdit      = widgets.NewQLineEdit(nil)
		aboutButton       = widgets.NewQPushButton(nil)

		sysTray = widgets.NewQSystemTrayIcon(nil)
	)

	sysTray.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))

	systrayMenu := widgets.NewQMenu(nil)
	openWinAction := systrayMenu.AddAction("打开界面")
	closeWinAction := systrayMenu.AddAction("关闭程序")

	sysTray.Show()

	expandSizePol := widgets.NewQSizePolicy()
	expandSizePol.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)

	StartButton.SetSizePolicy(expandSizePol)
	StopButton.SetSizePolicy(expandSizePol)
	aboutButton.SetSizePolicy(expandSizePol)

	StartButton.SetText("开始")
	StopButton.SetText("结束")
	aboutButton.SetText("关于")

	StopButton.SetEnabled(false)
	infoLineEdit.SetReadOnly(true)

	validatorLineEdit.SetValidator(gui.NewQIntValidator2(0, 65535, validatorLineEdit))
	validatorLineEdit.SetText("33899")

	infoLineEdit.SetSizePolicy(expandSizePol)
	// infoLineEdit.SetEnabled(false)

	StartButton.ConnectClicked(func(bool) {

		sserver.Start(validatorLineEdit.Text())
		infoLineEdit.SetText("start: " + validatorLineEdit.Text())
		isRunning := sserver.GetRunState()
		StartButton.SetEnabled(!isRunning)
		StopButton.SetEnabled(isRunning)
		validatorLineEdit.SetEnabled(!isRunning)

		go func() {
			// defer func() {
			// 	e := recover()
			// 	if e != nil {
			// 		infoLineEdit.SetText(fmt.Sprintf("%v", e))
			// 	}
			// }()

			errs := sserver.GetErrorOfType(socks5.ListenError{}, 3*time.Second)
			if len(errs) != 0 {
				isRunning := sserver.GetRunState()
				StartButton.SetEnabled(!isRunning)
				StopButton.SetEnabled(isRunning)
				validatorLineEdit.SetEnabled(!isRunning)
				for idx, listenErr := range errs {
					if listenErr == nil {
						continue
					}
					isRunning := sserver.GetRunState()
					StartButton.SetEnabled(!isRunning)
					StopButton.SetEnabled(isRunning)
					validatorLineEdit.SetEnabled(!isRunning)

					infoLineEdit.SetText(fmt.Sprintf("Listen err[%d]: %v", idx, listenErr))
				}

			}
		}()

	})
	StopButton.ConnectClicked(func(bool) {
		sserver.Stop()
		isRunning := sserver.GetRunState()
		StartButton.SetEnabled(!isRunning)
		StopButton.SetEnabled(isRunning)
		validatorLineEdit.SetEnabled(!isRunning)
		infoLineEdit.SetText(fmt.Sprintf("stop[%v]: %v", isRunning, validatorLineEdit.Text()))
	})

	aboutButton.ConnectClicked(func(checked bool) {
		widgets.QMessageBox_About(nil, "about qproxy", fmt.Sprintf("qproxy %s(%s)\nQt %s\n%s", verTag, verCommitHash, core.QtGlobal_qVersion(), strings.Replace(runtime.Version(), "go", "Go", -1)))
	})

	var validatorLayout = widgets.NewQGridLayout2()
	validatorLayout.AddWidget2(validatorLabel, 0, 0, 0)
	validatorLayout.AddWidget2(validatorLineEdit, 0, 1, 0)
	validatorLayout.AddWidget2(StartButton, 1, 0, 0)
	validatorLayout.AddWidget2(StopButton, 1, 1, 0)
	validatorLayout.AddWidget2(infoLineEdit, 2, 0, 0)
	validatorLayout.AddWidget2(aboutButton, 2, 1, 0)
	validatorGroup.SetLayout(validatorLayout)

	var layout = widgets.NewQGridLayout2()

	layout.AddWidget2(validatorGroup, 1, 0, 0)

	var window = NewHideWindow(nil, 0)
	window.SetWindowTitle("代理")

	openWinAction.ConnectTriggered(func(checked bool) {
		window.Show()
	})
	closeWinAction.ConnectTriggered(func(bool) {
		gui.QGuiApplication_SetQuitOnLastWindowClosed(true)
		sysTray.Hide()
		app.Quit()
	})
	sysTray.SetContextMenu(systrayMenu)
	sysTray.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		switch reason {
		case widgets.QSystemTrayIcon__DoubleClick:
			window.Show()
		}
	})

	var centralWidget = widgets.NewQWidget(window, 0)
	centralWidget.SetLayout(layout)

	window.SetCentralWidget(centralWidget)
	window.SetFixedHeight(150)
	window.SetFixedWidth(250)
	window.SetWindowIcon(gui.NewQIcon5(":/qml/icon.ico"))
	window.Show()
	gui.QGuiApplication_SetQuitOnLastWindowClosed(false)
	app.Exec()
}
