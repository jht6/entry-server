package entity

import "time"

// 灰度规则
type Rule struct {
	RuleId     uint      `json:"id" gorm:"primaryKey;column:id"`
	RuleName   string    `json:"rule_name" gorm:"column:rule_name"`
	RuleType   int       `json:"rule_type" gorm:"column:rule_type"`
	RuleStatus int       `json:"rule_status" gorm:"column:rule_status"`
	Percent    int       `json:"percent" gorm:"column:percent"`
	Entry      string    `json:"entry" gorm:"column:entry"`
	CreateUser string    `json:"creater" gorm:"column:creater"`
	UpdateUser string    `json:"updater" gorm:"column:updater"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	PublishId  int       `json:"publish_id" gorm:"column:publish_id"`
}

func (Rule) TableName() string {
	return "publish_rules"
}
