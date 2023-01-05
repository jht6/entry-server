package middleware

import (
	"encoding/json"
	"entry-server/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogReqRes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 记录
		logger := utils.CtxGetLogger(ctx)
		body, _ := json.Marshal(ctx.Request.Body)
		logger.WithFields(logrus.Fields{
			"host":        ctx.Request.Host,
			"rtx":         ctx.Request.Header.Get("staffname"),
			"staffid":     ctx.Request.Header.Get("staffid"),
			"path":        ctx.Request.URL.Path,
			"http_method": ctx.Request.Method,
			"http_query":  ctx.Request.URL.Query,
			"http_body":   string(body),
			"user_agent":  ctx.Request.Header.Get("user-agent"),
			"referer":     ctx.Request.Header.Get("referer"),
		}).Info("Request comes.")

		ctx.Next()

		// 记录响应状态码
		logger.WithFields(logrus.Fields{
			"status_code": ctx.Writer.Status(),
		}).Info("Response backs.")

	}
}
