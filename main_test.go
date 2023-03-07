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
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

// 前置处理
func TestMain(m *testing.M) {
	os.Setenv("SERVER_ENV", "test")
	router = setupRouter()

	// start another server to serve html
	cmd := exec.Command("bash", "-c", "go run . &")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Fail to start backend go server, please check!")
		return
	}

	// clear db
	db := utils.GetDB()
	db.Exec("delete from t_publish")
	db.Exec("delete from t_rule")

	// clear redis
	redis.FlushDB()

	m.Run()
}

func doPost(body *gin.H, path string) {
	bodyByte, _ := json.Marshal(*body)
	req := httptest.NewRequest("POST", path, bytes.NewReader(bodyByte))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
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
var testRuleName = "percent_test"

func TestApiCreateAndUpdateRule(t *testing.T) {
	testName := testRuleName
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
	ret := db.First(&rule, "name = ?", testName)
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

	//  ------ test update rule
	ruleId := rule.RuleId
	newStatus := 1
	body = gin.H{
		"rule_id":        ruleId,
		"status":         newStatus,
		"publish_domain": "tomcat",
	}
	jsonByte, _ = json.Marshal(body)
	req = httptest.NewRequest("POST", "/api/update_rule", bytes.NewReader(jsonByte))
	w = httptest.NewRecorder()
	req.Header.Set("content-type", "application/json")
	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &res)

	// check response
	assert.Equal(t, 200, w.Code, "Http code should be 200")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, int(res["code"].(float64)))
	assert.Equal(t, testName, res["data"].(map[string]any)["name"])
	assert.Equal(t, testDomain, res["data"].(map[string]any)["publish_domain"])
	assert.Equal(t, testEntry, res["data"].(map[string]any)["entry"])
	assert.Equal(t, newStatus, int(res["data"].(map[string]any)["status"].(float64)))

	// check db
	ret = db.First(&rule, "name = ?", testName)
	assert.Equal(t, int64(1), ret.RowsAffected)
	assert.Equal(t, testDomain, rule.PublishDomain) // 预期不被更新
	assert.Equal(t, newStatus, rule.Status)

	// check redis
	cachedData, err = redis.GetRuleListByDomain(testDomain)
	assert.Equal(t, nil, err)
	var y []map[string]any
	json.Unmarshal([]byte(cachedData), &y)
	assert.Equal(t, testName, y[0]["name"])
	assert.Equal(t, testType, int(y[0]["type"].(float64)))
	assert.Equal(t, testConfig, y[0]["config"])
	assert.Equal(t, testDomain, y[0]["publish_domain"])
	assert.Equal(t, testEntry, y[0]["entry"])
	assert.Equal(t, newStatus, int(y[0]["status"].(float64)))
}

// ====== end api:rule

// ====== start entry-server

func TestEntryServer(t *testing.T) {
	// case 1: 503 if publish doesn't exist
	t.Run("503 if publish doesn't exist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Host = "none.xx.com"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 503, w.Code)
	})

	// case 2: get main version if gray rules don't exist
	t.Run("main version", func(t *testing.T) {
		// prepare a publish without any rule
		domain := "test-main.es.com"
		body := gin.H{
			"name":   "test_main_version",
			"domain": domain,
			"entry":  "http://localhost:8080/html/main.html",
		}
		bodyByte, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/create_publish", bytes.NewReader(bodyByte))
		req.Header.Set("content-type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// expect to get main version html
		req = httptest.NewRequest("GET", "/", nil)
		req.Host = domain
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		html := w.Body.String()

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "<div>main</div>", html)
	})

	// case 3: gray rule
	t.Run("gray version", func(t *testing.T) {
		// prepare a publish and three rules
		domain := "test-gray.es.com"

		body := gin.H{
			"name":   "test_gray_rule",
			"domain": domain,
			"entry":  "http://localhost:8080/html/main.html",
		}
		doPost(&body, "/api/create_publish")

		// rule 1: somebody
		body = gin.H{
			"name":           "test_rule_somebody",
			"type":           2,
			"config":         "{\"user_list\":[1,2]}",
			"entry":          "http://localhost:8080/html/somebody.html",
			"publish_domain": domain,
		}
		doPost(&body, "/api/create_rule")

		// rule 2: header
		body = gin.H{
			"name":           "test_rule_header",
			"type":           3,
			"config":         "{\"header\":\"gray||||true\"}",
			"entry":          "http://localhost:8080/html/header.html",
			"publish_domain": domain,
		}
		doPost(&body, "/api/create_rule")

		// rule 3: percent
		body = gin.H{
			"name":           "test_rule_percent",
			"type":           1,
			"config":         "{\"percent\":10}",
			"entry":          "http://localhost:8080/html/percent.html",
			"publish_domain": domain,
		}
		doPost(&body, "/api/create_rule")

		t.Run("match somebody rule", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Host = domain
			req.Header.Set("user-id", "1")
			req.Header.Set("gray", "true")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			html := w.Body.String()

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, "<div>somebody</div>", html)
		})

		t.Run("match header rule", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Host = domain
			req.Header.Set("user-id", "20")
			req.Header.Set("gray", "true")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			html := w.Body.String()

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, "<div>header</div>", html)
		})

		t.Run("match percent rule", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Host = domain
			req.Header.Set("user-id", "9")
			req.Header.Set("gray", "false")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			html := w.Body.String()

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, "<div>percent</div>", html)
		})
	})
}

// ====== end entry-server
