package redis

import (
	"context"
	"strconv"

	"entry-server/common/utils"

	"github.com/go-redis/redis/v9"
)

const PREFIX_MAIN_ENTRY = "PUB:ENTRY:"
const PREFIX_HTMLURL_TO_HTMLCONTENT = "PUB:HTMLURL:TO:HTMLCONTENT:"
const PREFIX_HOST_TO_RULES = "PUB:RULES:"
const PREFIX_HOST_TO_PUBLISH = "ES:HOST:TO:PUBLISH:"

var ctx = context.Background()

func getConn() *redis.Client {
	addr := utils.CfgGet("REDIS_ADDR")
	pwd := utils.CfgGet("REDIS_PWD")
	db, _ := strconv.Atoi(utils.CfgGet("REDIS_DB"))

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
}

func FlushDB() *redis.StatusCmd {
	rdb := getConn()
	return rdb.FlushDB(ctx)
}

// publish config
func GetPublishByDomain(domain string) (string, error) {
	rdb := getConn()
	publish, err := rdb.Get(ctx, PREFIX_HOST_TO_PUBLISH+domain).Result()
	return publish, err
}

func SetPublish(domain string, publish string) error {
	rdb := getConn()
	err := rdb.Set(ctx, PREFIX_HOST_TO_PUBLISH+domain, publish, 0).Err()
	return err
}

// 将来源host与响应html文件url的对应关系存入redis
func SetHtmlUrl(host string, htmlUrl string) error {
	rdb := getConn()
	err := rdb.Set(ctx, PREFIX_MAIN_ENTRY+host, htmlUrl, 0).Err()
	return err
}

// 根据host读取html文件url
func GetHtmlUrlByHost(host string) (string, error) {
	rdb := getConn()
	url, err := rdb.Get(ctx, PREFIX_MAIN_ENTRY+host).Result()
	return url, err
}

// 缓存html内容
func SetHtmlContent(url string, content string) error {
	md5 := utils.StringToMd5(url)
	rdb := getConn()
	err := rdb.Set(ctx, PREFIX_HTMLURL_TO_HTMLCONTENT+md5, content, 0).Err()
	return err
}

// 根据htmlUrl读取html内容缓存
func GetHtmlContentByUrl(url string) (string, error) {
	md5 := utils.StringToMd5(url)
	rdb := getConn()
	html, err := rdb.Get(ctx, PREFIX_HTMLURL_TO_HTMLCONTENT+md5).Result()
	return html, err
}

// 根据host读取rules
func GetRulesByHost(host string) ([]string, error) {
	rdb := getConn()
	vals, err := rdb.ZRange(ctx, PREFIX_HOST_TO_RULES+host, 0, -1).Result()
	return vals, err
}
