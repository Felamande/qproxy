package main

import (
	"gopkg.in/ini.v1"
)

func main() {
	launchType := "qml"
	launchTypeMap := map[string]func(*ini.File){
		"qtwidget": QtMain,
		"qml":      QmlMain,
	}
	iniCfg, err := ini.Load("launch.ini")
	if err == nil {
		typ := iniCfg.Section("").Key("launch").MustString("qtwidget")
		if _, ok := launchTypeMap[typ]; ok {
			launchType = typ
		}
	}

	launchTypeMap[launchType](iniCfg)
}
