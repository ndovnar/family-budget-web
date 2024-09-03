package authz

import (
	"github.com/gin-gonic/gin"
)

func (a *Authz) IsUserHasReadAcessToTransaction(ctx *gin.Context, id string) bool {
	transaction, err := a.store.GetTransaction(ctx, id)
	if err != nil {
		return false
	}

	return a.IsUserHasAccessToCategory(ctx, transaction.CategoryID)
}

func (a *Authz) IsUserHasWriteAcessToTransaction(ctx *gin.Context, id string) bool {
	claims := a.auth.GetClaimsFromContext(ctx)

	transaction, err := a.store.GetTransaction(ctx, id)
	if err != nil {
		return false
	}

	return transaction.UserID == claims.ID
}
