package main

import (
	"os"

	"github.com/Felamande/qproxy/socks5server"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"gopkg.in/ini.v1"
)

func init() {

}

func QmlMain(ini *ini.File) {
	gui.NewQGuiApplication(len(os.Args), os.Args)
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	debug := ini.Section("").Key("qmldebug").MustBool(false)
	debugQmlFile := ini.Section("").Key("qmldebug_file").MustString("../../qml/view.qml")

	socks5server.Socks5Server_QmlRegisterType2("socks5", 1, 0, "socks5server")

	var app = qml.NewQQmlApplicationEngine(nil)

	if debug {
		app.Load(core.QUrl_FromLocalFile(debugQmlFile))
	} else {
		app.Load(core.NewQUrl3("qrc:/qml/view.qml", 0))

	}

	gui.QGuiApplication_Exec()

}
