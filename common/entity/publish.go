package entity

import "time"

// 发布项目
type Publish struct {
	PublishId  uint      `json:"publish_id" gorm:"primaryKey;column:publish_id"`
	Domain     string    `json:"domain" gorm:"column:domain"`
	Name       string    `json:"name" gorm:"column:name"`
	Entry      string    `json:"entry" gorm:"column:entry"`
	Status     uint      `json:"status" gorm:"column:status"`
	CreateUser string    `json:"create_user" gorm:"column:create_user"`
	UpdateUser string    `json:"update_user" gorm:"column:update_user"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Publish) TableName() string {
	return "t_publish"
}
