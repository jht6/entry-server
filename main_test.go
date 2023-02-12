package main

import (
	"bytes"
	"encoding/json"
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

}

func TestApiUpdatePublish(t *testing.T) {
	// TODO
}

// ====== api:publish end
