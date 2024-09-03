package authz

import "github.com/gin-gonic/gin"

func (a *Authz) IsUserHasAccessToBudget(ctx *gin.Context, id string) bool {
	claims := a.auth.GetClaimsFromContext(ctx)

	budget, err := a.store.GetBudget(ctx, id)
	if err != nil {
		return false
	}

	return budget.OwnerID == claims.UserID
}
