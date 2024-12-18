package authz

import "github.com/gin-gonic/gin"

func (a *Authz) IsUserHasAccessToUser(ctx *gin.Context, id string) bool {
	claims := a.auth.GetClaimsFromContext(ctx)

	return claims.UserID == id
}
