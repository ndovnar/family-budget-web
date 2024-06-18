package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, err := range ctx.Errors {
			switch e := err.Err.(type) {
			case error.HTTPError:
				if e.Description != "" {
					ctx.AbortWithStatusJSON(e.StatusCode, e)
				} else {
					ctx.AbortWithStatus(e.StatusCode)
				}
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
