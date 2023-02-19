package main

import (
	"bytes"
	"encoding/json"
	"entry-server/common/entity"
	"entry-server/common/redis"
	"entry-server/common/utils"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

// 前置处理
func TestMain(m *testing.M) {
	os.Setenv("SERVER_ENV", "test")
	router = setupRouter()

	// clear db
	db := utils.GetDB()
	db.Exec("delete from t_publish")

	// clear redis
	redis.FlushDB()

	m.Run()
}

func TestGet503(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	req.Host = "none.xx.com"

	router.ServeHTTP(w, req)

	assert.Equal(t, 503, w.Code, "Get 503 when the requested publish doesn't exist")
}

// ====== start api:publish
var testDomain string = "test1.es.com"

func TestApiCreatePublish(t *testing.T) {
	testName := "test1"
	// testDomain := "test1.es.com"
	testEntry := "http://localhost:8080/html/b.html"
	var testStatus int = 0
	body := gin.H{
		"name":   testName,
		"domain": testDomain,
		"entry":  testEntry,
	}
	jsonByte, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/create_publish", bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()

	req.Header.Set("content-type", "application/json")

	router.ServeHTTP(w, req)

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)

	// check response
	assert.Equal(t, 200, w.Code, "Http code should be 200")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, int(res["code"].(float64)))
	assert.Equal(t, testName, res["data"].(map[string]any)["name"])
	assert.Equal(t, testDomain, res["data"].(map[string]any)["domain"])
	assert.Equal(t, testEntry, res["data"].(map[string]any)["entry"])
	assert.Equal(t, testStatus, int(res["data"].(map[string]any)["status"].(float64)))

	// check db
	db := utils.GetDB()
	var publish entity.Publish
	ret := db.First(&publish, "domain = ?", testDomain)
	assert.Equal(t, int64(1), ret.RowsAffected)
	assert.Equal(t, testName, publish.Name)
	assert.Equal(t, testDomain, publish.Domain)
	assert.Equal(t, testEntry, publish.Entry)
	assert.Equal(t, testStatus, publish.Status)

	// check redis
	cachedData, err := redis.GetPublishByDomain(testDomain)
	assert.Equal(t, nil, err)
	var x map[string]any
	json.Unmarshal([]byte(cachedData), &x)
	assert.Equal(t, testName, x["name"])
	assert.Equal(t, testDomain, x["domain"])
	assert.Equal(t, testEntry, x["entry"])
	assert.Equal(t, testStatus, int(x["status"].(float64)))
}

func TestApiUpdatePublish(t *testing.T) {
	testName := "test1_update"
	testEntry := "http://localhost:8080/html/c.html"
	testStatus := 1

	body := gin.H{
		"domain": testDomain,
		"name":   testName,
		"entry":  testEntry,
		"status": testStatus,
	}
	jsonByte, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/update_publish", bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()

	req.Header.Set("content-type", "application/json")

	router.ServeHTTP(w, req)

	var res map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &res)

	// check response
	assert.Equal(t, 200, w.Code, "Http code should be 200")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, int(res["code"].(float64)))
	assert.Equal(t, testName, res["data"].(map[string]any)["name"])
	assert.Equal(t, testDomain, res["data"].(map[string]any)["domain"])
	assert.Equal(t, testEntry, res["data"].(map[string]any)["entry"])
	assert.Equal(t, testStatus, int(res["data"].(map[string]any)["status"].(float64)))

	// check db
	db := utils.GetDB()
	var publish entity.Publish
	ret := db.First(&publish, "domain = ?", testDomain)
	assert.Equal(t, int64(1), ret.RowsAffected)
	assert.Equal(t, testName, publish.Name)
	assert.Equal(t, testDomain, publish.Domain)
	assert.Equal(t, testEntry, publish.Entry)
	assert.Equal(t, testStatus, publish.Status)

	// check redis
	cachedData, err := redis.GetPublishByDomain(testDomain)
	fmt.Printf("redis: %v\n", cachedData)
	assert.Equal(t, nil, err)
	var x map[string]any
	json.Unmarshal([]byte(cachedData), &x)
	assert.Equal(t, testName, x["name"])
	assert.Equal(t, testDomain, x["domain"])
	assert.Equal(t, testEntry, x["entry"])
	assert.Equal(t, testStatus, int(x["status"].(float64)))
}

// ====== end api:publish

// ====== start api:rule

func TestApiCreateRule(t *testing.T) {
	testName := "percent"
	testType := 1
	testConfig := `{"percent":20}`
	testEntry := "http://localhost:8080/html/percent.html"
	testStatus := 0

	body := gin.H{
		"name":           testName,
		"type":           testType,
		"config":         testConfig,
		"entry":          testEntry,
		"publish_domain": testDomain,
	}
	jsonByte, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/create_rule", bytes.NewReader(jsonByte))
	w := httptest.NewRecorder()

	req.Header.Set("content-type", "application/json")

	router.ServeHTTP(w, req)

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)

	// check response
	assert.Equal(t, 200, w.Code, "Http code should be 200")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, int(res["code"].(float64)))
	assert.Equal(t, testName, res["data"].(map[string]any)["name"])
	assert.Equal(t, testType, int(res["data"].(map[string]any)["type"].(float64)))
	assert.Equal(t, testConfig, res["data"].(map[string]any)["config"])
	assert.Equal(t, testStatus, int(res["data"].(map[string]any)["status"].(float64)))
	assert.Equal(t, testEntry, res["data"].(map[string]any)["entry"])
	assert.Equal(t, testDomain, res["data"].(map[string]any)["publish_domain"])

	// check db
	db := utils.GetDB()
	var rule entity.Rule
	ret := db.First(&rule, "publish_domain = ?", testDomain)
	assert.Equal(t, int64(1), ret.RowsAffected)
	assert.Equal(t, testName, rule.Name)
	assert.Equal(t, testType, rule.Type)
	assert.Equal(t, testConfig, rule.Config)
	assert.Equal(t, testDomain, rule.PublishDomain)
	assert.Equal(t, testEntry, rule.Entry)
	assert.Equal(t, testStatus, rule.Status)

	// check redis
	cachedData, err := redis.GetRuleListByDomain(testDomain)
	assert.Equal(t, nil, err)
	var x []map[string]any
	json.Unmarshal([]byte(cachedData), &x)
	assert.Equal(t, testName, x[0]["name"])
	assert.Equal(t, testType, int(x[0]["type"].(float64)))
	assert.Equal(t, testConfig, x[0]["config"])
	assert.Equal(t, testDomain, x[0]["publish_domain"])
	assert.Equal(t, testEntry, x[0]["entry"])
	assert.Equal(t, testStatus, int(x[0]["status"].(float64)))
}

// ====== end api:rule
