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

// ====== api:publish start

func TestApiCreatePublish(t *testing.T) {
	testName := "test1"
	testDomain := "test1.es.com"
	testEntry := "http://localhost:8080/html/b.html"
	var testStatus uint = 0
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

	assert.Equal(t, 200, w.Code, "Http code should be 200")

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, int(res["code"].(float64)))
	assert.Equal(t, testName, res["data"].(map[string]any)["name"])
	assert.Equal(t, testDomain, res["data"].(map[string]any)["domain"])
	assert.Equal(t, testEntry, res["data"].(map[string]any)["entry"])
	assert.Equal(t, testStatus, uint(res["data"].(map[string]any)["status"].(float64)))

	// 校验db数据
	db := utils.GetDB()
	var publish entity.Publish
	ret := db.First(&publish, "domain = ?", testDomain)
	assert.Equal(t, int64(1), ret.RowsAffected)
	assert.Equal(t, testName, publish.Name)
	assert.Equal(t, testDomain, publish.Domain)
	assert.Equal(t, testEntry, publish.Entry)
	assert.Equal(t, testStatus, publish.Status)

	// 校验redis数据
	cachedData, err := redis.GetPublishByDomain(testDomain)
	fmt.Printf("rr: %v\n", cachedData)
	assert.Equal(t, nil, err)
	var x map[string]any
	json.Unmarshal([]byte(cachedData), &x)
	assert.Equal(t, testName, x["name"])
	assert.Equal(t, testDomain, x["domain"])
	assert.Equal(t, testEntry, x["entry"])
	assert.Equal(t, testStatus, uint(x["status"].(float64)))
}

func TestApiUpdatePublish(t *testing.T) {
	// TODO
	// 允许修改 domain name entry status

}

// ====== api:publish end
