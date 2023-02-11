package publish

import "time"

type CreatePublishDto struct {
	Name       string    `json:"name" binding:"required"`
	Domain     string    `json:"domain" binding:"required"`
	Entry      string    `json:"entry" binding:"required"`
	Status     uint      `json:"status"`
	CreateUser string    `json:"create_user"`
	UpdateUser string    `json:"update_user"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdatePublishDto struct {
	PublishId uint   `json:"publish_id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Entry     string `json:"entry"`
	Status    uint   `json:"status"`
}