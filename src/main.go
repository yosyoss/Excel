package main

import (
	get "Yosyos/src/api/getlist"
	save "Yosyos/src/api/save"
	A "Yosyos/src/api/test"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	v1 := router.Group("/api/v1")
	v1.Use(A.PassAuthMiddleware())
	{
		v1.GET("/common/excel/list", get.GetListExcelValue)
		v1.GET("/common/excel/save", save.ProcessData)
	}

	router.Run(":9696")
}
