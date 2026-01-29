package responses

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(OK, gin.H{
		"status":    "success",
		"timestamp": time.Now().Local(),
		"data":      data,
	})
}

func Created(ctx *gin.Context, data interface{}) {
	ctx.JSON(CREATED, gin.H{
		"status":    "resource created",
		"timestamp": time.Now().Local(),
		"data":      data,
	})
}

func BadRequest(ctx *gin.Context, err error) {
	ctx.JSON(BAD_REQUEST, gin.H{
		"timestamp": time.Now().Local(),
		"status":    "invalid request body",
		"error":     err.Error(),
	})
}

func InternalServerError(ctx *gin.Context, err error) {
	ctx.JSON(INTERNAL_SERVER_ERROR, gin.H{
		"timestamp": time.Now().Local(),
		"status":    "internal server error",
		"error":     err.Error(),
	})
}

func Unauthorized(ctx *gin.Context, err error) {
	ctx.JSON(UNAUTHORIZED, gin.H{
		"timestamp": time.Now().Local(),
		"status":    "unauthorized",
		"error":     err.Error(),
	})
}
