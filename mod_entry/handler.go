package mod_entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"entry-server/common/constant"
	"entry-server/common/dao"
	"entry-server/common/entity"
	"entry-server/common/redis"
	"entry-server/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func EntryHandler(ctx *gin.Context) {
	uuid := ctx.GetString("uuid")
	ctx.Header("X-Req-Uuid", uuid)
	ctx.Header("X-Served-By", "Entry-Server")

	logger := utils.CtxGetLogger(ctx)
	logger.Info("请求进入核心处理器 EntryHandler()")

	domain := utils.CtxGetHost(ctx)

	// 获取发布项配置
	publish := dao.GetPublishByDomain(domain)
	if publish == nil {
		logger.WithFields(logrus.Fields{
			"domain": domain,
		}).Error("未找到该domain对应的发布项目")

		// TODO 做一个好看的404页面
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	logger.WithFields(logrus.Fields{
		"domain": domain,
	}).Info("已找到domain对应的发布项,继续处理")

	// 先看是否命中灰度
	entry := getGrayEntry(publish, ctx, logger)
	if entry == "" {
		entry = publish.Entry
		logger.WithFields(logrus.Fields{
			"html_url": entry,
		}).Info("未命中任何灰度规则，读取默认entry")
	}

	// 尝试redis缓存
	html, err := redis.GetHtmlContentByUrl(entry)
	if err == nil {
		logger.Info("从redis读到html内容,直接响应")
		ctx.Data(http.StatusOK, "text/html", []byte(html))
		return
	}

	// 缓存未命中，请求html源文件
	logger.WithFields(logrus.Fields{
		"entry": entry,
	}).Info("redis中没有html内容,准备直接请求html源文件")
	res, err := http.Get(entry)
	if err != nil || res.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("请求html源文件出错: %v", err))
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	reader := res.Body
	defer reader.Close()

	html, err = utils.GetStringFromReader(reader)
	if err != nil {
		logger.Error(fmt.Sprintf("reader转string异常:%v", err))
	}

	// 缓存html
	redis.SetHtmlContent(entry, html)

	logger.Info("从源头获取html内容成功并响应")
	ctx.Data(http.StatusOK, "text/html", []byte(html))
}

// 根据灰度规则获取entry, 按以下顺序进行匹配
// 1. 指定用户规则
// 2. 指定header规则
// 3. 百分比规则
func getGrayEntry(publish *entity.Publish, ctx *gin.Context, logger *logrus.Entry) string {
	logger.Info("尝试从redis读取灰度规则")

	domain := publish.Domain
	rulesStr, err := redis.GetRuleListByDomain(domain)

	if err != nil {
		return ""
	}

	var rules []entity.Rule
	json.Unmarshal([]byte(rulesStr), &rules)

	userId, err := strconv.Atoi(ctx.Request.Header.Get("user-id"))
	if err != nil {
		logger.Info("将header中user-id转为int时出现异常")
		return ""
	}

	// 尝试匹配【指定用户】规则
	for _, rule := range rules {
		if rule.Type == constant.GRAY_RULE_TYPE_USER {
			// rule.Config json字符串 {users:["id1","id2"]}
			var config entity.RuleConfig
			json.Unmarshal([]byte(rule.Config), &config)

			for _, id := range config.UserList {
				if id == userId {
					// 匹配成功
					logger.WithFields(logrus.Fields{
						"rule_id":   rule.RuleId,
						"rule_name": rule.Name,
						"config":    rule.Config,
						"entry":     rule.Entry,
						"user_id":   userId,
					}).Info("命中指定用户灰度规则")
					return rule.Entry
				}
			}
		}
	}

	// 尝试匹配指定header
	for _, rule := range rules {
		if rule.Type == constant.GRAY_RULE_TYPE_HEADER {
			var config entity.RuleConfig
			json.Unmarshal([]byte(rule.Config), &config)

			// header配置形如 "x-gray||||true"
			tmp := strings.Split(config.Header, "||||")
			headerKey := tmp[0]
			headerValue := tmp[1]

			if ctx.Request.Header.Get(headerKey) == headerValue {
				// 匹配成功
				return rule.Entry
			}
		}
	}

	// 尝试匹配百分比规则
	for _, rule := range rules {
		if rule.Type == constant.GRAY_RULE_TYPE_PERCENT {
			var config entity.RuleConfig
			json.Unmarshal([]byte(rule.Config), &config)

			percent := config.Percent
			if isHitPercentRule(userId, percent) {
				return rule.Entry
			}
		}
	}

	return ""
}

// 员工id是否命中百分比
func isHitPercentRule(id int, percent int) bool {
	tail := id % 100
	return tail < percent
}
