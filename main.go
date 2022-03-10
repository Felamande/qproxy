package main

import (
	"github.com/spf13/viper"
)

func main() {
	launchType := defaultCfg["launch"].(string)
	launchTypeMap := map[string]func(){
		"qtwidget": QtMain,
		"qml":      QmlMain,
	}
	typ := viper.GetString("launch")
	if _, ok := launchTypeMap[typ]; ok {
		launchType = typ
	}

	launchTypeMap[launchType]()
}
