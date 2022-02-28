package main

import (
	"fmt"
	"os"

	"github.com/Felamande/go-socks5"
	"github.com/Felamande/qproxy/socks5server"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

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
		validatorGroup    = widgets.NewQGroupBox2("socks5代理", nil)
		validatorLabel    = widgets.NewQLabel2("端口:", nil, 0)
		validatorLineEdit = widgets.NewQLineEdit(nil)
		StartButton       = widgets.NewQPushButton(nil)
		StopButton        = widgets.NewQPushButton(nil)
		infoLineEdit      = widgets.NewQLineEdit(nil)

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
	StartButton.SetText("开始")
	StopButton.SetText("结束")
	StopButton.SetEnabled(false)
	infoLineEdit.SetReadOnly(true)

	validatorLineEdit.SetValidator(gui.NewQIntValidator2(0, 65535, validatorLineEdit))
	validatorLineEdit.SetText("33899")

	infoLineEdit.SetSizePolicy(expandSizePol)
	// infoLineEdit.SetEnabled(false)

	StartButton.ConnectClicked(func(bool) {
		sserver.Start(validatorLineEdit.Text())
		go func() {
			infoLineEdit.SetText("start: " + validatorLineEdit.Text())
			// StartButton.SetEnabled(false)
			StartButton.SetEnabled(false)
			StopButton.SetEnabled(true)
			validatorLineEdit.SetEnabled(false)

			errs := sserver.GetAllError(3)
			for _, err := range errs {
				switch errListenErr := err.(type) {
				case socks5.ListenError:
					infoLineEdit.SetText(fmt.Sprintf("Listen err: %v", errListenErr))
					StartButton.SetEnabled(true)
					StopButton.SetEnabled(false)
					validatorLineEdit.SetEnabled(true)
					return
				default:
					infoLineEdit.SetText(fmt.Sprintf("get other error: %v", errListenErr))
				}

			}

			StartButton.SetEnabled(false)
			StopButton.SetEnabled(true)
			validatorLineEdit.SetEnabled(false)

		}()

	})
	StopButton.ConnectClicked(func(bool) {
		sserver.Stop()
		StartButton.SetEnabled(true)
		StopButton.SetEnabled(false)
		validatorLineEdit.SetEnabled(true)
		infoLineEdit.SetText("stop: " + validatorLineEdit.Text())
	})

	var validatorLayout = widgets.NewQGridLayout2()
	validatorLayout.AddWidget2(validatorLabel, 0, 0, 0)
	validatorLayout.AddWidget2(validatorLineEdit, 0, 1, 0)
	validatorLayout.AddWidget2(StartButton, 1, 0, 0)
	validatorLayout.AddWidget2(StopButton, 1, 1, 0)
	validatorLayout.AddWidget2(infoLineEdit, 2, 0, 0)
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
