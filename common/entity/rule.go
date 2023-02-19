package entity

import "time"

// 灰度规则
type Rule struct {
	RuleId        uint      `json:"rule_id" gorm:"primaryKey;column:rule_id"`
	Name          string    `json:"name" gorm:"column:name"`
	Type          uint      `json:"type" gorm:"column:type"`
	Config        string    `json:"config" gorm:"column:config"`
	Status        uint      `json:"status" gorm:"column:status"`
	Entry         string    `json:"entry" gorm:"column:entry"`
	CreateUser    string    `json:"create_user" gorm:"column:create_user"`
	UpdateUser    string    `json:"update_user" gorm:"column:update_user"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
	PublishDomain string    `json:"publish_domain" gorm:"column:publish_domain"`
}

// 灰度规则配置结构
type RuleConfig struct {
	UserList []int  `json:"user_list"`
	Percent  int    `json:"percent"`
	Header   string `json:"header"`
}

func (Rule) TableName() string {
	return "t_rule"
}
