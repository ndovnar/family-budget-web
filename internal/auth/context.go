package auth

import "github.com/gin-gonic/gin"

const claimsContextKey = "authorization_payload"

func (*Auth) SetClaimsToContext(ctx *gin.Context, claims *Claims) {
	ctx.Set(claimsContextKey, claims)
}

func (*Auth) GetClaimsFromContext(ctx *gin.Context) *Claims {
	return ctx.MustGet(claimsContextKey).(*Claims)
}
