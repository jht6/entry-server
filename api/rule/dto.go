package rule

type CreateRuleDto struct {
	Name          string `json:"name" binding:"required"`
	Type          uint   `json:"type" binding:"required"`
	Config        string `json:"config" binding:"required"`
	Entry         string `json:"entry" binding:"required"`
	Status        uint   `json:"status"`
	CreateUser    string `json:"create_user"`
	PublishDomain string `json:"publish_domain" binding:"required"`
}