package entity

// 配置服务那边会自动将mysql中的字段转为驼峰，存入redis的也是驼峰
// 用下面结构体来反序列化redis中缓存的规则数据，将其处理成和DB实体一样的字段名

type CachedRule struct {
	Entry      string `json:"entry"`
	RuleId     uint   `json:"id"`
	Percent    int    `json:"percent"`
	PublishId  int    `json:"publishId"`
	RuleName   string `json:"ruleName"`
	RuleStatus int    `json:"ruleStatus"`
	RuleType   int    `json:"ruleType"`

	RuleUsers []CachedUser
	RuleOrgs  []interface{} `json:"ruleOrgs"`
}

type CachedUser struct {
	UserId      uint   `json:"userId"`
	EnglishName string `json:"englishName"`
}
