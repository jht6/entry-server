package main

import (
	"fmt"
	"strings"

	"entry-server/common/middleware"
	"entry-server/common/utils"
	"entry-server/mod_api"
	"entry-server/mod_entry"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitConfig()

	router := gin.Default()

	// 为了安全
	router.SetTrustedProxies(nil)

	router.Use(middleware.Uuid())
	router.Use(middleware.LogReqRes())

	router.Static("/html", "./html")

	router.GET("/test", func(ctx *gin.Context) {
		host := strings.Split(ctx.Request.Host, ":")[0]
		fmt.Printf("host:%v\n", host)
	})

	// 核心handler
	router.GET("/", mod_entry.EntryHandler)

	// 挂载api处理器
	mountApiHandler(router)

	router.Run(":8080")
}

func mountApiHandler(router *gin.Engine) {
	router.GET("/api/add", mod_api.AddCfgHandler)
	// router.POST("/api/create_project", mod_api.CreateProject)
}
