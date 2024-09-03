package authz

import (
	"github.com/gin-gonic/gin"
)

func (a *Authz) IsUserHasAccessToAccount(ctx *gin.Context, id string) bool {
	claims := a.auth.GetClaimsFromContext(ctx)

	account, err := a.store.GetAccount(ctx, id)
	if err != nil {
		return false
	}

	return account.OwnerID == claims.UserID
}
