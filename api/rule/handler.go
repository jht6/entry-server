package rule

import (
	"encoding/json"
	"entry-server/common/entity"
	"entry-server/common/redis"
	"entry-server/common/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetRuleList(ctx *gin.Context) {
	// 根据域名查询相关规则
	var dto GetRuleDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	db := utils.GetDB()
	var rules []entity.Rule
	ret := db.Where("publish_domain = ?", dto.PublishDomain).Find(&rules)
	if ret.Error != nil {
		utils.CtxResAbort(ctx, ret.Error.Error())
		return
	}

	utils.CtxResOk(ctx, rules)
}

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

	domain := dto.PublishDomain
	rule := entity.Rule{
		Name:          dto.Name,
		Type:          dto.Type,
		Config:        dto.Config,
		Status:        dto.Status,
		Entry:         dto.Entry,
		PublishDomain: domain,
	}

	ret = db.Create(&rule)
	if ret.Error != nil {
		utils.CtxResAbort(ctx, ret.Error.Error())
		return
	}

	// 将domain关联的所有灰度规则重新缓存，包括未启用的规则
	var rules []entity.Rule
	db.Where("publish_domain = ?", domain).Find(&rules)
	byteRules, _ := json.Marshal(rules)
	redis.SetRuleListByDomain(domain, string(byteRules))

	utils.CtxResOk(ctx, rule)
}

func UpdateRuleHandler(ctx *gin.Context) {
	var dto map[string]any
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		utils.CtxResAbort(ctx, err.Error())
		return
	}

	ruleId := int(dto["rule_id"].(float64))
	db := utils.GetDB()
	var rule entity.Rule
	ret := db.First(&rule, "rule_id = ?", ruleId)
	if ret.RowsAffected == 0 {
		utils.CtxResAbort(ctx, fmt.Sprintf("不存在rule_id=[%v]的规则", ruleId))
		return
	}

	domain := rule.PublishDomain

	// 更新相关字段配置
	fields := map[string]int{
		"name":        0,
		"type":        0,
		"config":      0,
		"entry":       0,
		"status":      0,
		"update_user": 0,
	}
	updatedData := map[string]any{}
	for k, v := range dto {
		if _, ok := fields[k]; ok {
			updatedData[k] = v
		}
	}
	db.Model(&rule).Updates(updatedData)

	// 将domain关联的所有灰度规则重新缓存，包括未启用的规则
	var rules []entity.Rule
	db.Where("publish_domain = ?", domain).Find(&rules)
	byteRules, _ := json.Marshal(rules)
	redis.SetRuleListByDomain(domain, string(byteRules))

	utils.CtxResOk(ctx, rule)

}
