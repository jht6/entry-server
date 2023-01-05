package entity

import "time"

// 发布项目
type Publish struct {
	PublishId  uint      `json:"id" gorm:"primaryKey;column:id"`
	Domain     string    `json:"domain" gorm:"column:domain"`
	Name       string    `json:"name" gorm:"column:name"`
	Entry      string    `json:"entry" gorm:"column:entry"`
	Status     uint      `json:"status" gorm:"column:status"`
	CreateUser string    `json:"creater" gorm:"column:creater"`
	UpdateUser string    `json:"updater" gorm:"column:updater"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Publish) TableName() string {
	return "publishes"
}
