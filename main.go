package main

import (
	"fmt"
	"strings"

	apiPublish "entry-server/api/publish"
	"entry-server/common/middleware"
	"entry-server/common/utils"
	"entry-server/mod_entry"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
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

	return router
}

func mountApiHandler(router *gin.Engine) {
	router.POST("/api/create_publish", apiPublish.CreatePublishHandler)
	router.POST("/api/update_publish", apiPublish.UpdatePublishHandler)
}
