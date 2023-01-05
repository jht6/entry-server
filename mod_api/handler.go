package mod_api

import (
	"fmt"

	"entry-server/common/redis"

	"github.com/gin-gonic/gin"
)

func AddCfgHandler(ctx *gin.Context) {
	fmt.Printf("add_cfg\n")

	err := redis.SetHtmlUrl("jht1.woa.com", "http://localhost:8080/html/a.html")
	if err != nil {
		panic(err)
	}

	fmt.Printf("set redis OK.\n")
}

// func CreateProject(ctx *gin.Context) {
// 	var json CreateProjectDto
// 	if err := ctx.ShouldBindJSON(&json); err != nil {
// 		utils.CtxResAbort(ctx, err.Error())
// 		return
// 	}

// 	db := utils.GetDB()

// 	// 检查是否已存在该host
// 	var project Project
// 	ret := db.First(&project, "host = ?", json.Host)
// 	if ret.RowsAffected != 0 {
// 		utils.CtxResAbort(ctx, fmt.Sprintf("域名 [%s] 已存在，请更换其他域名", json.Host))
// 		return
// 	}

// 	project = Project{
// 		ProjectName: json.ProjectName,
// 		Host:        json.Host,
// 		HtmlUrl:     json.HtmlUrl,
// 		CreateUser:  0, // TODO 鉴权
// 		UpdateUser:  0, // TODO 鉴权
// 	}

// 	ret = db.Create(&project)
// 	if ret.Error != nil {
// 		utils.CtxResAbort(ctx, ret.Error.Error())
// 		return
// 	}

// 	// 同步redis
// 	redis.SetHtmlUrl(project.Host, project.HtmlUrl)

// 	utils.CtxResOk(ctx, project)
// }
