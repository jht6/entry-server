package entity

import "time"

// 灰度用户
type RuleOrg struct {
	RuleOrgId     uint      `json:"id" gorm:"primaryKey;column:id"`
	OrgId         uint      `json:"org_id" gorm:"column:org_id"`
	Name          string    `json:"name" gorm:"column:name"`
	FullName      string    `json:"full_name" gorm:"column:full_name"`
	CreateUser    string    `json:"creater" gorm:"column:creater"`
	UpdateUser    string    `json:"updater" gorm:"column:updater"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
	PublishRuleId uint      `json:"publish_rule_id" gorm:"column:publish_rule_id"`
}

func (RuleOrg) TableName() string {
	return "publish_rule_orgs"
}
