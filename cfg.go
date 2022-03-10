package main

import "github.com/spf13/viper"

var defaultCfg = map[string]interface{}{
	"launch":        "qml",
	"qmldebug":      false,
	"qmldebug_file": "./qml/view.qml",
}

func init() {
	viper.SetConfigName("launch")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")

	viper.ReadInConfig()
	if defaultSecITF, exist := viper.AllSettings()["default"]; exist {
		switch defaultSec := defaultSecITF.(type) {
		case map[string]interface{}:
			for k, v := range defaultSec {
				viper.Set(k, v)
			}
		}
	}
	for k, v := range defaultCfg {
		viper.SetDefault(k, v)
	}
}
