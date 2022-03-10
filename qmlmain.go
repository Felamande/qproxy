package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/Felamande/qproxy/socks5server"
	"github.com/spf13/viper"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
)

type VerGetter struct {
	core.QObject

	_ func() `constructor:"init"`
	_ string `property:"VerCommitHash"`
	_ string `property:"VerTag"`
	_ string `property:"GoVersion"`
	_ string `property:"QtVersion"`
}

func (v *VerGetter) init() {
	v.SetVerCommitHash(verCommitHash)
	v.SetVerTag(verTag)
	gover := runtime.Version()
	gover = strings.Replace(gover, "go", "", -1)
	v.SetGoVersion(gover)
	v.SetQtVersion(core.QtGlobal_qVersion())
}

func QmlMain() {

	var debug bool
	var debugQmlFile string

	gapp := gui.NewQGuiApplication(len(os.Args), os.Args)
	gapp.SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	debug = viper.GetBool("qmldebug")
	debugQmlFile = viper.GetString("qmldebug_file")

	socks5server.Socks5Server_QmlRegisterType2("Socks5", 1, 0, "Socks5server")

	var app = qml.NewQQmlApplicationEngine(nil)
	app.RootContext().SetContextProperty("VerGetter", NewVerGetter(nil))

	if debug {
		if f, err := os.Stat(debugQmlFile); f == nil || err != nil {
			app.Load(core.NewQUrl3("qrc:/qml/view.qml", 0))
		} else {
			app.Load(core.QUrl_FromLocalFile(debugQmlFile))
		}
	} else {
		app.Load(core.NewQUrl3("qrc:/qml/view.qml", 0))
	}

	gapp.SetQuitOnLastWindowClosed(false)
	gapp.Exec()

}
