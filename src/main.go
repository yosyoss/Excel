package main

import (
	ctrl "Yosyos/src/controller"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("src/config/Config.json")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}
}
func main() {
	ctrl.RunController()
	return
}
