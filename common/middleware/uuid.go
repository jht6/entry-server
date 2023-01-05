package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Uuid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New().String()
		ctx.Set("uuid", id)
		ctx.Next()
	}
}
