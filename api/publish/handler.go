package publish

import (
	"encoding/json"
	"entry-server/common/utils"
	"fmt"

	"entry-server/common/entity"
	"entry-server/common/redis"

	"github.com/gin-gonic/gin"
)

func CreatePublishHandler(ctx *gin.Context) {
	var dto CreatePublishDto
	dto.Status = 0 // 默认启用
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	db := utils.GetDB()

	var publish entity.Publish

	// 检查是否已存在该host配置
	ret := db.First(&publish, "domain = ?", dto.Domain)
	if ret.RowsAffected != 0 {
		utils.CtxResAbort(ctx, fmt.Sprintf("域名 [%s] 已存在，创建失败", dto.Domain))
		return
	}

	publish = entity.Publish{
		Domain: dto.Domain,
		Name:   dto.Name,
		Entry:  dto.Entry,
		Status: dto.Status,
	}

	ret = db.Create(&publish)
	if ret.Error != nil {
		utils.CtxResAbort(ctx, ret.Error.Error())
		return
	}

	// 同步redis
	jsonbyte, _ := json.Marshal(publish)
	redis.SetPublish(dto.Domain, string(jsonbyte))

	utils.CtxResOk(ctx, publish)
}

func UpdatePublishHandler(ctx *gin.Context) {
	// 查询publish_id的项
	var dto UpdatePublishDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	db := utils.GetDB()

	var publish entity.Publish

	ret := db.First(&publish, "publish_id = ?", dto.PublishId)
	if ret.RowsAffected == 0 {
		utils.CtxResAbort(ctx, fmt.Sprintf("不存在ID=%v的发布项", dto.PublishId))
		return
	}

	// 若查到了,更新相关字段配置
	db.Model(&publish).Updates(dto)

	// 同步redis
	jsonbyte, _ := json.Marshal(publish)
	redis.SetPublish(dto.Domain, string(jsonbyte))

	utils.CtxResOk(ctx, publish)
}
