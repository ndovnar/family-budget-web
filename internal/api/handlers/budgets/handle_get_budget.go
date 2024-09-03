package budgets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ndovnar/family-budget-api/internal/api/error"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/rs/zerolog/log"
)

func (b *Budgets) HandleGetBudget(ctx *gin.Context) {
	id := ctx.Param("id")

	hasAccess := b.authz.IsUserHasAccessToBudget(ctx, id)
	if !hasAccess {
		ctx.Error(error.NewHttpError(http.StatusForbidden))
		return
	}

	budget, err := b.store.GetBudget(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get budget")

		if err == store.ErrNotFound {
			ctx.Error(error.NewHttpError(http.StatusNotFound))
		} else {
			ctx.Error(error.NewHttpError(http.StatusInternalServerError))
		}

		return
	}

	ctx.JSON(http.StatusOK, budget)
}
