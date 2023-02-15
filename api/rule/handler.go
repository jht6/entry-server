package rule

import (
	"entry-server/common/entity"
	"entry-server/common/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateRuleHandler(ctx *gin.Context) {
	var dto CreateRuleDto
	dto.Status = 0
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	// TODO 检查config是否为合法json

	db := utils.GetDB()
	var publish entity.Publish
	// 检查domain对应的publish是否存在
	ret := db.First(&publish, "domain = ?", dto.PublishDomain)
	if ret.RowsAffected == 0 {
		utils.CtxResAbort(ctx, fmt.Sprintf("未找到域名 [%s] 对应的发布项目，无法创建灰度规则", dto.PublishDomain))
		return
	}

	rule := entity.Rule{
		Name:          dto.Name,
		Type:          dto.Type,
		Config:        dto.Config,
		Status:        dto.Status,
		Entry:         dto.Entry,
		PublishDomain: dto.PublishDomain,
	}

	ret = db.Create(&rule)
	if ret.Error != nil {
		utils.CtxResAbort(ctx, ret.Error.Error())
		return
	}

	// TODO 缓存到redis

	utils.CtxResOk(ctx, rule)
}
