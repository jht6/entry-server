package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CtxGetHost(ctx *gin.Context) string {
	return strings.Split(ctx.Request.Host, ":")[0]
}

func CtxResOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}

func CtxResOkNil(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
}

func CtxResAbort(ctx *gin.Context, msg string) {
	errCode := 400
	CtxResAbortWithCode(ctx, msg, errCode)
}

func CtxResAbortWithCode(ctx *gin.Context, msg string, code int) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func CtxGetLogger(ctx *gin.Context) *logrus.Entry {
	uuid := ctx.GetString("uuid")
	logger := logrus.New()
	if IsProd() {
		logger.Formatter = &logrus.JSONFormatter{}

		file, err := os.OpenFile("/data/log/entry-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.Out = file
		} else {
			fmt.Printf("打开entry-server.log文件失败\n")
		}
	}

	entry := logger.WithFields(logrus.Fields{
		"uuid": uuid,
	})
	return entry
}
