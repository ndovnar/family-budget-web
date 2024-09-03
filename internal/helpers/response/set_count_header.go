package response

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetCountHeader(ctx *gin.Context, count int64) {
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
}
