package publish

import (
	"entry-server/common/utils"
	"fmt"

	"entry-server/common/entity"

	"github.com/gin-gonic/gin"
)

func CreatePublishHandler(ctx *gin.Context) {
	var json CreatePublishDto
	json.Status = 0 // 默认启用
	if err := ctx.ShouldBindJSON(&json); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	db := utils.GetDB()

	var publish entity.Publish

	// 检查是否已存在该host配置
	ret := db.First(&publish, "domain = ?", json.Domain)
	if ret.RowsAffected != 0 {
		utils.CtxResAbort(ctx, fmt.Sprintf("域名 [%s] 已存在，创建失败", json.Domain))
		return
	}

	publish = entity.Publish{
		Domain: json.Domain,
		Name:   json.Name,
		Entry:  json.Entry,
		Status: json.Status,
	}

	ret = db.Create(&publish)
	if ret.Error != nil {
		utils.CtxResAbort(ctx, ret.Error.Error())
		return
	}

	// TODO 同步redis

	utils.CtxResOk(ctx, publish)
}
