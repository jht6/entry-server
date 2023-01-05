package entity

import "time"

// 灰度用户
type RuleUser struct {
	RuleUserId    uint      `json:"id" gorm:"primaryKey;column:id"`
	UserId        uint      `json:"user_id" gorm:"column:user_id"`
	ChineseName   string    `json:"chinese_name" gorm:"column:chinese_name"`
	EnglishName   string    `json:"english_name" gorm:"column:english_name"`
	FullName      string    `json:"full_name" gorm:"column:full_name"`
	CreateUser    string    `json:"creater" gorm:"column:creater"`
	UpdateUser    string    `json:"updater" gorm:"column:updater"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
	PublishRuleId uint      `json:"publish_rule_id" gorm:"column:publish_rule_id"`
}

func (RuleUser) TableName() string {
	return "publish_rule_users"
}
