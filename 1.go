package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Felamande/go-socks5"
	"github.com/Felamande/qproxy/socks5server"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var verTag string
var verCommitHash string

type Socks5ServerGroupBox struct {
	widgets.QGroupBox

	portLabel     *widgets.QLabel
	portLineInput *widgets.QLineEdit
	startButton   *widgets.QPushButton
	stopButton    *widgets.QPushButton
	infoLineEdit  *widgets.QLineEdit

	s5s *socks5server.Socks5Server

	layout *widgets.QGridLayout
}

func (gb *Socks5ServerGroupBox) Init() *Socks5ServerGroupBox {

	gb.SetTitle("socks5代理")
	gb.portLabel = widgets.NewQLabel2("端口:", nil, 0)
	gb.portLineInput = widgets.NewQLineEdit2("33899", nil)
	gb.startButton = widgets.NewQPushButton2("开始", nil)
	gb.stopButton = widgets.NewQPushButton2("结束", nil)
	gb.infoLineEdit = widgets.NewQLineEdit(nil)

	expandSizePol := widgets.NewQSizePolicy()
	expandSizePol.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
	gb.startButton.SetSizePolicy(expandSizePol)
	gb.stopButton.SetSizePolicy(expandSizePol)

	gb.stopButton.SetEnabled(false)
	gb.infoLineEdit.SetReadOnly(true)

	gb.portLineInput.SetValidator(gui.NewQIntValidator2(0, 65535, gb.portLineInput))

	gb.layout = widgets.NewQGridLayout2()
	gb.layout.AddWidget2(gb.portLabel, 0, 0, 0)
	gb.layout.AddWidget2(gb.portLineInput, 0, 1, 0)
	gb.layout.AddWidget2(gb.startButton, 1, 0, 0)
	gb.layout.AddWidget2(gb.stopButton, 1, 1, 0)
	gb.layout.AddWidget2(gb.infoLineEdit, 2, 0, 0)
	gb.SetLayout(gb.layout)

	gb.s5s = socks5server.NewSocks5Server(nil).Init()

	gb.startButton.ConnectClicked(func(checked bool) {
		gb.s5s.StartServer(gb.portLineInput.Text())
	})
	gb.stopButton.ConnectClicked(func(checked bool) {
		gb.s5s.StopServer()
	})

	gb.s5s.ConnectRunStateChange(gb.runStateChange)
	gb.s5s.ConnectReceiveRunningError(gb.processRunningError)

	return gb

}

func (gb *Socks5ServerGroupBox) processRunningError(iErr interface{}) {

	switch err := iErr.(type) {
	case socks5.ListenError:
		gb.infoLineEdit.SetText(err.Error())
	default:
		return
	}
}

func (gb *Socks5ServerGroupBox) runStateChange(isRunning bool) {
	gb.startButton.SetEnabled(!isRunning)
	gb.stopButton.SetEnabled(isRunning)
	gb.portLineInput.SetEnabled(!isRunning)
}

type ProxyAppWidget struct {
	widgets.QWidget

	s5sGroupBox *Socks5ServerGroupBox
	layout      *widgets.QGridLayout
}

func (w *ProxyAppWidget) Init() *ProxyAppWidget {
	w.s5sGroupBox = NewSocks5ServerGroupBox(nil).Init()

	w.layout = widgets.NewQGridLayout(nil)
	w.layout.AddWidget2(w.s5sGroupBox, 1, 0, 0)
	w.SetLayout(w.layout)

	return w
}

type ProxyAppTray struct {
	widgets.QSystemTrayIcon

	menu    *widgets.QMenu
	actions map[string]*widgets.QAction
}

func (t *ProxyAppTray) Init() *ProxyAppTray {
	t.actions = make(map[string]*widgets.QAction)
	t.menu = widgets.NewQMenu(nil)
	t.SetContextMenu(t.menu)

	return t
}

func (t *ProxyAppTray) AddAction(text string, triggerFn func(bool)) *ProxyAppTray {
	action := t.menu.AddAction(text)

	action.ConnectTriggered(triggerFn)
	t.actions[text] = action

	return t
}

type ProxyAppWindow struct {
	widgets.QMainWindow

	appWidget *ProxyAppWidget

	tray *ProxyAppTray
}

func (w *ProxyAppWindow) Init() *ProxyAppWindow {
	w.appWidget = NewProxyAppWidget(w, 0).Init()
	w.SetCentralWidget(w.appWidget)

	w.tray = NewProxyAppTray(nil).Init()
	w.tray.SetIcon(gui.NewQIcon5(":/qml/icon.ico"))

	w.tray.AddAction("打开界面", func(bool) {
		w.Show()
	}).AddAction("关闭程序", func(bool) {
		gui.QGuiApplication_SetQuitOnLastWindowClosed(true)
		w.tray.Hide()
		core.QCoreApplication_Instance().Quit()
	}).AddAction("关于程序", func(b bool) {
		widgets.QMessageBox_About(nil, "about qproxy", fmt.Sprintf("qproxy %s(%s)\nQt %s\n%s", verTag, verCommitHash, core.QtGlobal_qVersion(), strings.Replace(runtime.Version(), "go", "Go", -1)))
	})

	w.tray.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		switch reason {
		case widgets.QSystemTrayIcon__DoubleClick:
			w.Show()
		}
	})

	return w
}

func (w *ProxyAppWindow) ShowTray() {
	w.tray.Show()
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := NewProxyAppWindow(nil, 0).Init()
	window.SetFixedHeight(150)
	window.SetFixedWidth(250)
	window.SetWindowIcon(gui.NewQIcon5(":/qml/icon.ico"))

	window.Show()
	window.ShowTray()

	gui.QGuiApplication_SetQuitOnLastWindowClosed(false)
	app.Exec()
}
