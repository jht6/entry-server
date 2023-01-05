package mod_api

import (
	"entry-server/common/entity"
	"entry-server/common/utils"
	"fmt"
)

func GetHtmlUrlByHost(host string) (string, error) {
	db := utils.GetDB()

	var publish entity.Publish

	ret := db.Where("domain = ?", host).First(&publish)

	if ret.RowsAffected == 0 {
		// 没查到
		utils.LogInfo(fmt.Sprintf("未查询到htmlUrl，host为%v", host))
		return "", ret.Error
	}

	return publish.Entry, nil
}
