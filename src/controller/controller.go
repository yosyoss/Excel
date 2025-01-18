package controller

import (
	get "Yosyos/src/api/getlist"
	save "Yosyos/src/api/save"
	A "Yosyos/src/api/test"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunController() {
	router := gin.New()

	v1 := router.Group("/api/v1")
	v1.Use(A.PassAuthMiddleware())
	{
		v1.GET("/common/excel/list", get.GetListExcelValue)
		v1.GET("/common/excel/save", save.ProcessData)
	}

	RunPort := viper.GetString("port")
	router.Run(RunPort)
}
