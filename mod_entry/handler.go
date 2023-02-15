package mod_entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	host := utils.CtxGetHost(ctx)

	// 获取发布项配置
	publish := dao.GetPublishByDomain(host)
	if publish == nil {
		logger.WithFields(logrus.Fields{
			"domain": host,
		}).Error("未找到该domain对应的发布项目")

		// TODO 做一个好看的404页面
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	logger.WithFields(logrus.Fields{
		"domain": host,
	}).Info("已找到domain对应的发布项,继续处理")

	// 先看是否命中灰度
	htmlUrl := getGrayEntry(publish, ctx, logger)
	if htmlUrl == "" {
		htmlUrl = publish.Entry
		logger.WithFields(logrus.Fields{
			"html_url": htmlUrl,
		}).Info("未命中任何灰度规则，读取默认htmlUrl")
	}

	// 尝试redis缓存
	html, err := redis.GetHtmlContentByUrl(htmlUrl)
	if err == nil {
		logger.Info("从redis读到html内容,直接响应")
		ctx.Data(http.StatusOK, "text/html", []byte(html))
		return
	}

	// 缓存未命中，请求html源文件
	logger.WithFields(logrus.Fields{
		"html_url": htmlUrl,
	}).Info("redis中没有html内容,准备直接请求html源文件")
	res, err := http.Get(htmlUrl)
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
	redis.SetHtmlContent(htmlUrl, html)

	logger.Info("从源头获取html内容成功并响应")
	ctx.Data(http.StatusOK, "text/html", []byte(html))
}

// 根据灰度规则获取htmlUrl, 按以下顺序进行匹配
// 1. 指定用户规则
// 2. 组织架构规则
// 3. 百分比规则
func getGrayEntry(publish *entity.Publish, ctx *gin.Context, logger *logrus.Entry) string {
	// 优先读取redis
	logger.Info("尝试从redis读取灰度规则")
	cachedRuleList, err := redis.GetRulesByHost(publish.Domain)

	if err != nil || len(cachedRuleList) == 0 {
		if err != nil {
			logger.Info("读取redis灰度规则出错，即将从mysql读取灰度规则")
		} else {
			logger.Info("未能从redis读取到灰度规则，即将从mysql读取灰度规则")
		}

		// 从mysql中读取
		rules := dao.GetRulesByPublishId(publish.PublishId)

		if rules == nil {
			return ""
		}

		url := getUserRuleUrl(rules, ctx, logger)
		if url != "" {
			return url
		}

		var percentRules []percentRule
		for _, v := range rules {

			if v.Type == constant.GRAY_RULE_TYPE_PERCENT {
				percentRules = append(percentRules, percentRule{
					RuleId:   v.RuleId,
					RuleName: v.Name,
					Percent:  10, // TODO
					Entry:    v.Entry,
					RuleType: v.Type,
				})
			}
		}
		url = getEntryByPercentRule(percentRules, ctx, logger)

		return url
	}

	logger.Info("读取redis灰度规则成功,开始匹配,规则条数:", len(cachedRuleList))

	var cachedUserRules []entity.CachedRule
	var cachedPercentRules []percentRule
	for _, str := range cachedRuleList {
		var rule entity.CachedRule
		err := json.Unmarshal([]byte(str), &rule)
		if err != nil {
			logger.Warn("反序列化规则json时出错, err=" + err.Error())
			continue
		}

		if rule.RuleType == constant.GRAY_RULE_TYPE_USER {
			cachedUserRules = append(cachedUserRules, rule)
		} else if rule.RuleType == constant.GRAY_RULE_TYPE_PERCENT {
			cachedPercentRules = append(cachedPercentRules, percentRule{
				RuleId:   rule.RuleId,
				RuleName: rule.RuleName,
				Percent:  10, // TODO 实现
				Entry:    rule.Entry,
				RuleType: uint(rule.RuleType),
			})
		}

		// TODO 补齐组织架构规则
	}

	// 指定用户
	entry := getEntryByCachedUserRule(cachedUserRules, ctx)
	if entry != "" {
		return entry
	}

	// 百分比
	entry = getEntryByPercentRule(cachedPercentRules, ctx, logger)

	return entry

}

func getEntryByCachedUserRule(rules []entity.CachedRule, ctx *gin.Context) string {
	username := ctx.Request.Header.Get("staffname")
	// 没有用户名，算作匹配失败，直接跳过
	if username == "" {
		return ""
	}

	for _, rule := range rules {
		for _, user := range rule.RuleUsers {
			if user.EnglishName == username {
				return rule.Entry
			}
		}
	}

	return ""
}

// 尝试匹配指定人员规则 for mysql
func getUserRuleUrl(rules []entity.Rule, ctx *gin.Context, logger *logrus.Entry) string {
	username := ctx.Request.Header.Get("staffname")
	// 没有用户名，算作匹配失败，直接跳过
	if username == "" {
		return ""
	}

	list := dao.GetUserIdsByRules(rules)

	for _, v := range list {
		if v.Type == constant.GRAY_RULE_TYPE_USER && v.EnglishName == username {
			logger.WithFields(logrus.Fields{
				"username":  username,
				"rule_id":   v.RuleId,
				"rule_type": v.Type,
			}).Info("命中指定用户规则")

			return v.Entry
		}
	}

	return ""
}

// 尝试匹配百分比规则
type percentRule struct {
	RuleId   uint
	RuleName string
	Percent  uint
	Entry    string
	RuleType uint
}

func getEntryByPercentRule(rules []percentRule, ctx *gin.Context, logger *logrus.Entry) string {
	rawStaffId := ctx.Request.Header.Get("staffid")
	staffId, err := strconv.Atoi(rawStaffId)

	// 若转换staffid失败，无需再继续处理
	if err != nil {
		logger.Warn(fmt.Sprintf("转换staffid时错误, err=%v", err))
		return ""
	}

	for _, v := range rules {
		if v.RuleType == constant.GRAY_RULE_TYPE_PERCENT && isHitPercentRule(staffId, int(v.Percent)) {
			logger.WithFields(logrus.Fields{
				"rule_id":   v.RuleId,
				"rule_name": v.RuleName,
				"percent":   v.Percent,
				"html_url":  v.Entry,
				"staff_id":  staffId,
			}).Info("命中百分比灰度规则")

			return v.Entry
		}
	}

	return ""
}

// 员工id是否命中百分比
func isHitPercentRule(id int, percent int) bool {
	tail := id % 100
	return tail < percent
}
