package dao

import (
	"encoding/json"
	"entry-server/common/constant"
	"entry-server/common/entity"
	"entry-server/common/redis"
	"entry-server/common/utils"
)

func GetPublishByDomain(domain string) *entity.Publish {
	var publish entity.Publish

	str, err := redis.GetPublishByDomain(domain)
	if err != nil {
		return nil
	}

	json.Unmarshal([]byte(str), &publish)

	return &publish
}

func GetRulesByPublishId(publishId uint) []entity.Rule {
	// TODO 读redis？

	db := utils.GetDB()

	var rules []entity.Rule
	enable := 1
	ret := db.Where("publish_id = ? AND rule_status = ?", publishId, enable).Find(&rules)

	if ret.RowsAffected == 0 {
		return nil
	}

	return rules
}

type RuleAndUser struct {
	entity.Rule
	UserId      uint   `gorm:"column:user_id"`
	EnglishName string `gorm:"column:english_name"`
}

func GetUserIdsByRules(rules []entity.Rule) []RuleAndUser {
	// 用户灰度规则id列表
	var ruleIds []uint
	for _, v := range rules {
		if v.Type == constant.GRAY_RULE_TYPE_USER {
			ruleIds = append(ruleIds, v.RuleId)
		}
	}

	db := utils.GetDB()

	var list []RuleAndUser
	db = db.Table("publish_rules")
	db = db.Select("publish_rules.*, publish_rule_users.user_id, publish_rule_users.english_name")
	db = db.Joins("left join publish_rule_users on publish_rule_users.publish_rule_id = publish_rules.id")
	db.Where("publish_rules.id in ?", ruleIds).Scan(&list)

	return list
}
