package entity

import "time"

// 灰度用户
type Token struct {
	TokenId     uint      `json:"id" gorm:"primaryKey;column:id"`
	Token       string    `json:"token" gorm:"column:token"`
	LandunToken string    `json:"landun_token" gorm:"column:landun_token"`
	CreateUser  string    `json:"creater" gorm:"column:creater"`
	UpdateUser  string    `json:"updater" gorm:"column:updater"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Token) TableName() string {
	return "publish_tokens"
}
