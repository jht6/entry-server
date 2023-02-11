package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestGet503(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	req.Host = "none.xx.com"

	router.ServeHTTP(w, req)

	assert.Equal(t, 503, w.Code, "Get 503 when the requested publish doesn't exist")
}

func TestMain(m *testing.M) {
	os.Setenv("SERVER_ENV", "test")
	router = setupRouter()
	m.Run()
}
