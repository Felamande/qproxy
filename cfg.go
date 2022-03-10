package main

import "github.com/spf13/viper"

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
	viper.SetDefault("launch", "qml")
	viper.SetDefault("qmldebug", false)
	viper.SetDefault("qmldebug_file", "./qml/view.qml")
}
